package api

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
)

// https://www.thepolyglotdeveloper.com/2020/02/interacting-with-a-graphql-api-with-golang/

var client = graphql.NewClient("http://localhost:3000/graphbrainz")

type ArtistInfoQuery struct {
	Search struct {
		Artists struct {
			Edges []struct {
				Node struct {
					MBID     string `json:"mbid"`
					Name     string `json:"name"`
					LifeSpan struct {
						Begin string `json:"begin"`
						End   string `json:"end"`
						Ended bool   `json:"ended"`
					} `json:"lifeSpan"`
					Country        string `json:"country"`
					Disambiguation string `json:"disambiguation"`
					ReleaseGroups  struct {
						Edges []struct {
							Node struct {
								Title           string `json:"title"`
								CoverArtArchive struct {
									Artwork bool   `json:"artwork"`
									Front   string `json:"front"`
								} `json:"coverArtArchive"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"releaseGroups"`
					Tags struct {
						Nodes []struct {
							Name  string `json:"name"`
							Count int    `json:"count"`
						} `json:"nodes"`
					} `json:"tags"`
					TheAudioDB struct {
						Biography string `json:"biography"`
					} `json:"theAudioDB"`
					Type string `json:"type"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"artists"`
	} `json:"search"`
}

func GetArtistInfo(name string) (ArtistInfoQuery, error) {
	request := graphql.NewRequest(`
	query GetArtistInfo($name: String!) {
		search {
			artists(query: $name, first: 1) {
				edges {
					node {
						mbid
						name
						lifeSpan {
							begin
							end
							ended
						}
						country
						disambiguation
						releaseGroups {
							edges {
								node {
									title
									coverArtArchive {
										artwork
										front
									}
								}
							}
						}
						tags {
							nodes {
								name
								count
							}
						}
						type
					}
				}
			}
		}
	}
	`)

	request.Var("name", name)

	var response ArtistInfoQuery
	if err := client.Run(context.Background(), request, &response); err != nil {
		return ArtistInfoQuery{}, fmt.Errorf("GetArtistInfo: %v", err)
	}

	return response, nil
}
