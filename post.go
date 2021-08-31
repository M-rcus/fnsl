package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Post struct {
	Author      string
	Id          string
	AccountId   string
	Content     string
	CreatedAt   time.Time
	Attachments []API_Common_Media
}

func (p *Post) DownloadAttachments(destDir string) {
	for _, att := range p.Attachments {
		if !att.IsOk() {
			log.Println("[Error] Skipping attachment not ok:", att.Id)
			continue
		}
		exists, err := att.LocalExists(destDir)
		if err != nil {
			log.Printf("[Error] Error checking for local file %s: %v\n", att.Id, err)
			continue
		}
		if exists {
			//log.Println("[Info] Attachment exists:", att.Id)
			continue
		}
		err = att.LocalDownload(destDir)
		if err != nil {
			log.Printf("[Error] Error downloading local file %s: %v\n", att.Id, err)
			continue
		}
	}
}

func (p *Post) LocalPostFilename() string {
	return fmt.Sprintf(
		"%s_%s.md",
		p.CreatedAt.UTC().Format("20060102-150405"),
		p.Id,
	)
}

func (p *Post) LocalPostExists(outDir string) (bool, error) {
	_, err := os.Stat(filepath.Join(outDir, p.LocalPostFilename()))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (p *Post) CreateLocalPost(postsDir, attachmentsDir string) error {
	f, err := os.Create(filepath.Join(postsDir, p.LocalPostFilename()))
	if err != nil {
		return err
	}
	_, err = f.WriteString(p.GetContent(postsDir, attachmentsDir))
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) GetContent(postsDir, attachmentsDir string) string {
	content := p.Content
	if p.Content != "" {
		content += "\n\n"
	}
	for _, att := range p.Attachments {
		if !att.IsOk() {
			continue
		}
		fn := filepath.Join(attachmentsDir, att.LocalFilename())
		rel, err := filepath.Rel(postsDir, fn)
		if err != nil {
			continue
		}

		switch {
		case att.IsPhoto():
			content += "![Photo](" + rel + ")\n\n"
		case att.IsVideo():
			attr := ""
			if att.Width > 0 {
				attr += fmt.Sprintf(`width="%d" `, att.Width)
			}
			if att.Height > 0 {
				attr += fmt.Sprintf(`height="%d" `, att.Height)
			}
			content += fmt.Sprintf(
				`<video controls %s style="max-width: 90%%"><source src="%s" type="%s"></video>`,
				attr,
				rel,
				att.MimeType,
			) + "\n\n"
		default:
			continue
		}
	}

	d := p.CreatedAt.UTC()
	return fmt.Sprintf(`---
title: "%s"
slug: "%s"
author: "%s"
date: "%s"
---

%s`,
		p.Author+" on "+d.Format(time.RFC822),
		p.Author+"/"+d.Format("20060102-150405")+"/"+p.Id,
		p.Author,
		d.Format(time.RFC3339),
		content,
	)
}
