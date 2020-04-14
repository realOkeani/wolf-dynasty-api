package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tmnt "github.homedepot.com/dev-insights/team-management-api"
	. "github.homedepot.com/dev-insights/team-management-api/http"
	"github.homedepot.com/dev-insights/team-management-api/models"
	"github.homedepot.com/dev-insights/team-management-api/sql/sqlfakes"
)

var _ = Describe("Teams", func() {

	var (
		svr           *httptest.Server
		uri           string
		teams         []models.Team
		t             time.Time
		router        *mux.Router
		fakeSQLClient *sqlfakes.FakeClient
	)

	BeforeEach(func() {
		router = mux.NewRouter().StrictSlash(true)
		fakeSQLClient = new(sqlfakes.FakeClient)
		s := tmnt.Services{
			SQLClient: fakeSQLClient,
		}
		AddRoutes(s, router)
		svr = httptest.NewServer(router)
		uri = svr.URL + "/v1/teams"
	})

	AfterEach(func() {
		svr.Close()
	})

	Describe("GET", func() {
		var (
			req *http.Request
			res *http.Response
			err error
		)

		BeforeEach(func() {
			req, _ = http.NewRequest("GET", uri, nil)
		})

		JustBeforeEach(func() {
			res, err = http.DefaultClient.Do(req)
		})

		Context("when getting the teams", func() {
			Context("and the server is down", func() {
				BeforeEach(func() {
					svr.Close()
				})

				It("returns an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("and the application is running", func() {
				Context("when the db returns an error", func() {
					BeforeEach(func() {
						fakeSQLClient.GetTeamsReturns([]models.Team{}, errors.New("some error"))
					})

					It("returns status internal server error", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})

				Context("when the db returns a list of teams", func() {
					BeforeEach(func() {
						t = time.Now()
						teams = []models.Team{
							{
								ID:        "1",
								Name:      "Awesome Team",
								CreatedAt: t,
								UpdatedAt: t,
							},
						}
						fakeSQLClient.GetTeamsReturns(teams, nil)
					})

					It("returns status OK", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(res.StatusCode).To(Equal(http.StatusOK))

						var retTeams []models.Team
						err = json.NewDecoder(res.Body).Decode(&retTeams)
						Expect(err).NotTo(HaveOccurred())

						Expect(retTeams[0].ID).To(Equal(teams[0].ID))
						Expect(retTeams[0].Name).To(Equal(teams[0].Name))
					})
				})
			})
		})
	})

	Describe("POST", func() {
		var (
			req  *http.Request
			res  *http.Response
			err  error
			team models.Team
			body *bytes.Buffer
		)

		BeforeEach(func() {
			body = &bytes.Buffer{}
		})

		JustBeforeEach(func() {
			req, _ = http.NewRequest("POST", uri, ioutil.NopCloser(body))
			res, err = http.DefaultClient.Do(req)
		})

		Context("when adding the team", func() {
			Context("and the server is down", func() {
				BeforeEach(func() {
					svr.Close()
				})

				It("returns an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("and the application is running", func() {
				Context("when the db returns an error", func() {
					BeforeEach(func() {
						t = time.Now()
						team = models.Team{
							ID:        "1",
							Name:      "Awesome Team",
							CreatedAt: t,
							UpdatedAt: t,
						}
						bodyBytes, err := json.Marshal(team)
						Expect(err).NotTo(HaveOccurred())

						body.Write(bodyBytes)
						fakeSQLClient.AddTeamReturns(models.Team{}, errors.New("some error"))
					})

					It("returns status internal server error", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})

				Context("when there is no body", func() {
					BeforeEach(func() {
						body.Write([]byte{})
					})

					It("returns status internal server error", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					})

				})

				Context("when a bad body is provided", func() {
					BeforeEach(func() {
						body.Write([]byte("{bad:json}"))
					})

					It("returns status internal server error", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					})

				})

				Context("when the db returns a list of teams", func() {
					BeforeEach(func() {
						t = time.Now()
						team = models.Team{
							ID:        "1",
							Name:      "Awesome Team",
							CreatedAt: t,
							UpdatedAt: t,
						}
						bodyBytes, err := json.Marshal(team)
						Expect(err).NotTo(HaveOccurred())

						body.Write(bodyBytes)
						fakeSQLClient.AddTeamReturns(team, nil)
					})

					It("returns status OK", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(res.StatusCode).To(Equal(http.StatusCreated))

						var retTeam models.Team
						err = json.NewDecoder(res.Body).Decode(&retTeam)
						Expect(err).NotTo(HaveOccurred())

						Expect(retTeam.ID).To(Equal(team.ID))
						Expect(retTeam.Name).To(Equal(team.Name))
					})
				})
			})
		})
	})
	Describe("PATCH", func() {
		var (
			req  *http.Request
			res  *http.Response
			err  error
			team models.Team
			body *bytes.Buffer
			guid string
		)

		BeforeEach(func() {
			guid = "some-guid"
			body = &bytes.Buffer{}
		})

		JustBeforeEach(func() {
			url := fmt.Sprintf("%s/%s", uri, guid)
			req, _ = http.NewRequest("PATCH", url, ioutil.NopCloser(body))
			res, err = http.DefaultClient.Do(req)
		})

		Context("when updating a team", func() {
			Context("and the server is down", func() {
				BeforeEach(func() {
					svr.Close()
				})

				It("returns an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("and the application is running", func() {
				Context("when there is no body", func() {
					BeforeEach(func() {
						body.Write([]byte{})
					})

					It("returns status internal server error", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					})
				})

				Context("when a bad body is provided", func() {
					BeforeEach(func() {
						body.Write([]byte("{bad:json}"))
					})

					It("returns status internal server error", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					})
				})

				Context("when the getting the team from the db", func() {
					BeforeEach(func() {
						t = time.Now()
						newTeam := models.Team{
							Name: "New Better Team",
						}
						bodyBytes, err := json.Marshal(newTeam)
						Expect(err).NotTo(HaveOccurred())

						body.Write(bodyBytes)
					})
					Context("returns the 'no rows found' error", func() {
						BeforeEach(func() {
							fakeSQLClient.GetTeamReturns(models.Team{}, errors.New("no rows found"))
						})

						It("returns status not found", func() {
							Expect(err).ToNot(HaveOccurred())
							Expect(res.StatusCode).To(Equal(http.StatusNotFound))
						})
					})

					Context("returns some error", func() {
						BeforeEach(func() {
							fakeSQLClient.GetTeamReturns(models.Team{}, errors.New("some error"))
						})

						It("returns status internal server error", func() {
							Expect(err).ToNot(HaveOccurred())
							Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
						})
					})

					Context("returns a team", func() {
						BeforeEach(func() {
							team = models.Team{
								ID:        "1",
								Name:      "Some Lame team",
								CreatedAt: t,
								UpdatedAt: t,
							}
							fakeSQLClient.GetTeamReturns(team, nil)
						})
						Context("when updating the team", func() {
							Context("returns an error", func() {
								BeforeEach(func() {
									fakeSQLClient.UpdateTeamReturns(models.Team{}, errors.New("some error"))
								})
								It("returns status internal server error", func() {
									Expect(err).ToNot(HaveOccurred())
									Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
								})
							})

							Context("is successful", func() {
								BeforeEach(func() {
									team = models.Team{
										ID:        "1",
										Name:      "Awesome Team",
										CreatedAt: t,
										UpdatedAt: t,
									}
									fakeSQLClient.UpdateTeamReturns(team, nil)
								})

								It("returns status OK", func() {
									Expect(err).ToNot(HaveOccurred())
									Expect(res.StatusCode).To(Equal(http.StatusOK))

									var retTeam models.Team
									err = json.NewDecoder(res.Body).Decode(&retTeam)
									Expect(err).NotTo(HaveOccurred())

									Expect(retTeam.ID).To(Equal(team.ID))
									Expect(retTeam.Name).To(Equal(team.Name))
								})
							})
						})
					})
				})
			})
		})
	})
	Describe("DELETE", func() {
		var (
			req  *http.Request
			res  *http.Response
			err  error
			team models.Team
			body *bytes.Buffer
			guid string
		)

		BeforeEach(func() {
			guid = "some-guid"
			body = &bytes.Buffer{}
		})

		JustBeforeEach(func() {
			url := fmt.Sprintf("%s/%s", uri, guid)
			req, _ = http.NewRequest("DELETE", url, ioutil.NopCloser(body))
			res, err = http.DefaultClient.Do(req)
		})

		Context("when deleting a team", func() {
			Context("and the server is down", func() {
				BeforeEach(func() {
					svr.Close()
				})

				It("returns an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("and the application is running", func() {
				Context("when getting the team from the db", func() {
					BeforeEach(func() {
						t = time.Now()
						newTeam := models.Team{
							Name: "New Better Team",
						}
						bodyBytes, err := json.Marshal(newTeam)
						Expect(err).NotTo(HaveOccurred())

						body.Write(bodyBytes)
					})
					Context("returns the 'no rows found' error", func() {
						BeforeEach(func() {
							fakeSQLClient.GetTeamReturns(models.Team{}, errors.New("no rows found"))
						})

						It("returns status not found", func() {
							Expect(err).ToNot(HaveOccurred())
							Expect(res.StatusCode).To(Equal(http.StatusNotFound))
						})
					})

					Context("returns some error", func() {
						BeforeEach(func() {
							fakeSQLClient.GetTeamReturns(models.Team{}, errors.New("some error"))
						})

						It("returns status internal server error", func() {
							Expect(err).ToNot(HaveOccurred())
							Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
						})
					})

					Context("returns a team", func() {
						BeforeEach(func() {
							team = models.Team{
								ID:        "1",
								Name:      "Some Lame team",
								CreatedAt: t,
								UpdatedAt: t,
							}
							fakeSQLClient.GetTeamReturns(team, nil)
						})
						Context("when deleting the team", func() {
							Context("returns an error", func() {
								BeforeEach(func() {
									fakeSQLClient.DeleteTeamReturns(errors.New("some error"))
								})
								It("returns status internal server error", func() {
									Expect(err).ToNot(HaveOccurred())
									Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
								})
							})

							Context("is successful", func() {
								BeforeEach(func() {
									team = models.Team{
										ID:        "1",
										Name:      "Awesome Team",
										CreatedAt: t,
										UpdatedAt: t,
									}
									fakeSQLClient.DeleteTeamReturns( nil)
								})

								It("returns status NoContent and no body", func() {
									Expect(err).ToNot(HaveOccurred())
									Expect(res.StatusCode).To(Equal(http.StatusNoContent))

									var retTeam models.Team
									err = json.NewDecoder(res.Body).Decode(&retTeam)
									Expect(err.Error()).To(Equal("EOF"))
								})
							})
						})
					})
				})
			})
		})
	})

})

