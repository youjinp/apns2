package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/youjinp/apns2"
	"github.com/youjinp/apns2/certificate"
	"github.com/youjinp/apns2/payload"
)

var (
	certificatePath string
	topic           string
	mode            string
	token           string
	title           string
	body            string
)

var rootCmd = &cobra.Command{
	Use:   "apns2",
	Short: "A tool for sending APNS notifications",
	Long:  `APNS2 is a CLI tool that helps with sending APNS notifications.`,
	Run: func(cmd *cobra.Command, args []string) {
		cert, pemErr := certificate.FromPemFile(certificatePath, "")

		if pemErr != nil {
			log.Fatalf("Error retrieving certificate `%v`: %v", certificatePath, pemErr)
		}

		client := apns2.NewClient(cert)

		if mode == "development" {
			client.Development()
		} else {
			client.Production()
		}

		res, err := client.Push(&apns2.Notification{
			DeviceToken: token,
			Topic:       topic,
			Payload:     payload.NewPayload().AlertTitle(title).AlertBody(body),
		})

		if err != nil {
			log.Fatal("Error: ", err)
		} else {
			fmt.Printf("%v: '%v'\n", res.StatusCode, res.Reason)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of APNS2",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("APNS2 v0.20.3")
	},
}

func main() {

	// root flags
	rootCmd.Flags().StringVarP(&certificatePath, "certificate-path", "c", "", "Path to certificate file.")
	rootCmd.Flags().StringVarP(&topic, "topic", "t", "", "The topic of the remote notification, which is typically the bundle ID for your app")
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "production", "APNS server to send notifications to. `production` or `development`")
	rootCmd.Flags().StringVarP(&token, "token", "d", "", "The device token to send notifications to")
	rootCmd.Flags().StringVarP(&title, "title", "e", "APNS Test", "The title of the APNS notification")
	rootCmd.Flags().StringVarP(&body, "body", "b", "APNS Test", "The body of the APNS notification")
	rootCmd.MarkFlagRequired("certificate-path")
	rootCmd.MarkFlagRequired("topic")
	rootCmd.MarkFlagRequired("token")

	// version
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
