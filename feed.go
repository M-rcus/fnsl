package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"
)

type Feed struct {
	Client *APIClient
}

func (x *Feed) GetPosts(profileId string, profileName string) (<-chan Post, error) {
	pch := make(chan Post)

	go func() {
		before := "0"
		for {
			profileUrl := fmt.Sprintf("https://apiv2.fansly.com/api/v1/timeline/%s?before=%s&after=0", profileId, before)
			read := &API_Timeline{}
			err := x.Client.GetAPI(profileUrl, read)
			if err != nil {
				log.Printf("[Error] Error requesting timeline (profile: %s, marker: %s): %v\n", profileId, before, err)
				break
			}

			if !read.Success {
				log.Printf("[Error] Timeline response is not successful (profile: %s, marker: %s)\n", profileId, before)
				break
			}

			if len(read.Response.Posts) == 0 {
				break
			}

			media := x.getMediaMap(&read.Response)
			mediaBundles := x.getMediaBundleMap(&read.Response)

			for _, p := range read.Response.Posts {
				if p.Id == "" || p.CreatedAt == 0 {
					continue
				}

				before = p.Id

				if len(p.Attachments) > 0 {
					sort.Slice(p.Attachments, func(i, j int) bool {
						posI, _ := strconv.Atoi(p.Attachments[i].Pos)
						posJ, _ := strconv.Atoi(p.Attachments[j].Pos)
						return posI < posJ
					})
				}

				obj := Post{
					Author:    profileName,
					Id:        p.Id,
					AccountId: p.AccountId,
					Content:   p.Content,
					CreatedAt: time.Unix(p.CreatedAt, 0),
				}

				attachments := make([]API_Common_Media, 0)
				for _, att := range p.Attachments {
					if att.ContentId == "" || att.ContentType == 0 || att.Pos == "" {
						continue
					}

					m, ok := media[att.ContentId]
					if ok && m.Media != nil {
						attachments = append(attachments, *m.Media)
					} else {
						m, ok := mediaBundles[att.ContentId]
						if ok && len(m.BundleContent) > 0 {
							for _, bc := range m.BundleContent {
								m, ok := media[bc.AccountMediaId]
								if ok && m.Media != nil {
									attachments = append(attachments, *m.Media)
								}
							}
						}
					}
				}
				if len(attachments) > 0 {
					obj.Attachments = attachments
				}

				pch <- obj
			}
		}

		close(pch)
	}()

	return pch, nil
}

func (x *Feed) getMediaMap(res *API_Timeline_Response) map[string]API_Timeline_Response_AccountMedia {
	mediaMap := make(map[string]API_Timeline_Response_AccountMedia, len(res.AccountMedia))
	for _, v := range res.AccountMedia {
		if v.Id == "" || !v.Access || v.CreatedAt < 1 {
			continue
		}
		mediaMap[v.Id] = v
	}
	return mediaMap
}

func (x *Feed) getMediaBundleMap(res *API_Timeline_Response) map[string]API_Timeline_Response_AccountMediaBundle {
	mediaBundleMap := make(map[string]API_Timeline_Response_AccountMediaBundle, len(res.AccountMediaBundles))
	for _, v := range res.AccountMediaBundles {
		if v.Id == "" || !v.Access || v.CreatedAt < 1 || len(v.AccountMediaIds) == 0 {
			continue
		}

		sort.Slice(v.BundleContent, func(i, j int) bool {
			return v.BundleContent[i].Pos < v.BundleContent[j].Pos
		})

		mediaBundleMap[v.Id] = v
	}
	return mediaBundleMap
}
