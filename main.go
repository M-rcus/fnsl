package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func main() {
	err := loadConfig()
	if err != nil {
		log.Fatalln("Error loading configuration:", err)
	}

	outBaseDir, err := filepath.Abs(viper.GetString("outDir"))
	if err != nil {
		log.Fatalln("Error determining output directory:", err)
	}

	api := &APIClient{
		Authorization: viper.GetString("authorization"),
		UserAgent:     viper.GetString("userAgent"),
	}
	api.Init()

	feedClient := &Feed{
		Client: api,
	}

	subsClient := &Subscriptions{
		Client: api,
	}

	subs, err := subsClient.GetSubscriptions()
	if err != nil {
		log.Fatalln("Failed to get list of subscriptions:", err)
	}

	if len(subs) == 0 {
		log.Fatalln("You're not subscribed to any profile")
	}

	for _, sub := range subs {
		log.Printf("Scraping profile %s (%s)\n", sub.Username, sub.Id)

		if !sub.ProfileAccess {
			log.Println("[Error] No access to profile", sub.Id)
			continue
		}

		outDir := filepath.Join(outBaseDir, sub.Username)
		postsDir := filepath.Join(outDir, "posts")
		attachmentsDir := filepath.Join(outDir, "attachments")
		err = os.MkdirAll(postsDir, 0777)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v\n", postsDir, err)
		}
		err = os.MkdirAll(attachmentsDir, 0777)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v\n", attachmentsDir, err)
		}

		pch, err := feedClient.GetPosts(sub.Id, sub.Username)
		if err != nil {
			log.Fatalf("Failed to get list of posts for profile %s: %v\n", sub.Id, err)
		}

		for post := range pch {
			post.DownloadAttachments(attachmentsDir)
			exists, err := post.LocalPostExists(postsDir)
			if err != nil {
				log.Println("[Error] Failed to check if local post exists:", post.Id)
				continue
			}
			if !exists {
				err = post.CreateLocalPost(postsDir, attachmentsDir)
				if err != nil {
					log.Println("[Error] Failed to write local post:", post.Id)
					continue
				}
			}
		}
	}
}

func loadConfig() error {
	viper.SetDefault("outDir", "out")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.fnsly")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	if viper.GetString("authorization") == "" {
		return errors.New("'authorization' missing or empty")
	}
	if viper.GetString("userAgent") == "" {
		return errors.New("'userAgent' missing or empty")
	}

	return nil
}
