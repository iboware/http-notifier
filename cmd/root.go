/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iboware/http-notifier/config"
	"github.com/iboware/http-notifier/pkg/notification"
	"github.com/spf13/cobra"
)

var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "http-notifier",
	Short: "Sends notifications to http APIs",
	Long: `It reads STDIN and sends new messages every interval. 
	Each line is interpreted as a new message that needs to be notified about. 
	It will keep running until it receives SIGINT, or OS Interrupt.`,
	Run: func(cmd *cobra.Command, args []string) {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, os.Interrupt)

		httpClient := &http.Client{
			Timeout: time.Second * 15,
		}
		logger := log.New(os.Stderr, "", log.LstdFlags)

		numberOfWorkers := 5
		channelSize := 3
		w, err := notification.NewNotificationPool(logger, numberOfWorkers, channelSize)
		if err != nil {
			panic(err)
		}
		w.Start()

		t := time.NewTicker(time.Second * time.Duration(cfg.Interval))
		scanner := bufio.NewScanner(os.Stdin)
	loop:
		for {
			select {
			case <-signals:
				fmt.Println("shutting down gracefully...")
				t.Stop()
				w.Stop()
				break loop
			case <-t.C:
				if scanner.Scan() {
					nt := notification.NewNotificationTask(cfg.Url, scanner.Text(), httpClient, func(err error) {
						log.Printf("notification could not sent: %t", err)
					})

					log.Printf("notification queued %v", nt)
					w.AddWorkNonBlocking(nt)
				}

				if err := scanner.Err(); err != nil {
					logger.Println(err)
				}
			}
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cfg.RegisterFlags(rootCmd.Flags())
	rootCmd.MarkFlagRequired("url")
}
