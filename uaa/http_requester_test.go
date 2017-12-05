package uaa_test

import (
	. "code.cloudfoundry.org/uaa-cli/uaa"

	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("HttpGetter", func() {
	var (
		server       *ghttp.Server
		client       *http.Client
		config       Config
		responseJson string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = &http.Client{}
		config = NewConfigWithServerURL(server.URL())
		responseJson = `{"foo": "bar"}`
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("UnauthenticatedRequester", func() {
		Describe("Get", func() {
			It("calls an endpoint with Accept application/json header", func() {
				server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
				))

				UnauthenticatedRequester{}.Get(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("supports zone switching", func() {
				server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))
				config.ZoneSubdomain = "twilight-zone"
				UnauthenticatedRequester{}.Get(client, config, "/testPath", "someQueryParam=true")
			})

			It("returns helpful error when GET request fails", func() {
				server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(500, ""),
					ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
				))

				_, err := UnauthenticatedRequester{}.Get(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})
		})

		Describe("Delete", func() {
			It("calls an endpoint with Accept application/json header", func() {
				server.RouteToHandler("DELETE", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("DELETE", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
				))

				UnauthenticatedRequester{}.Delete(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("supports zone switching", func() {
				server.RouteToHandler("DELETE", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))
				config.ZoneSubdomain = "twilight-zone"
				UnauthenticatedRequester{}.Delete(client, config, "/testPath", "someQueryParam=true")
			})

			It("returns helpful error when DELETE request fails", func() {
				server.RouteToHandler("DELETE", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(500, ""),
					ghttp.VerifyRequest("DELETE", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
				))

				_, err := UnauthenticatedRequester{}.Delete(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})
		})

		Describe("PostForm", func() {
			It("calls an endpoint with correct body and headers", func() {
				responseJson = `{
				  "access_token" : "bc4885d950854fed9a938e96b13ca519",
				  "token_type" : "bearer",
				  "expires_in" : 43199,
				  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
				  "jti" : "bc4885d950854fed9a938e96b13ca519"
				}`

				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
					ghttp.VerifyBody([]byte("hello=world")),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				))

				body := map[string]string{"hello": "world"}
				returnedBytes, _ := UnauthenticatedRequester{}.PostForm(client, config, "/oauth/token", "", body)
				parsedResponse := string(returnedBytes)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(parsedResponse).To(ContainSubstring("expires_in"))
			})

			It("treats 201 as success", func() {
				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(201, responseJson),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
				))

				_, err := UnauthenticatedRequester{}.PostForm(client, config, "/oauth/token", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).To(BeNil())
			})

			It("treats 405 as error", func() {
				server.RouteToHandler("PUT", "/oauth/token/foo/secret", ghttp.CombineHandlers(
					ghttp.RespondWith(405, responseJson),
					ghttp.VerifyRequest("PUT", "/oauth/token/foo/secret", ""),
				))

				_, err := UnauthenticatedRequester{}.PutJson(client, config, "/oauth/token/foo/secret", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
			})

			It("returns an error when request fails", func() {
				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(500, "garbage"),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
				))

				_, err := UnauthenticatedRequester{}.PostForm(client, config, "/oauth/token", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})

			It("supports zone switching", func() {
				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(201, responseJson),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
					ghttp.VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))

				config.ZoneSubdomain = "twilight-zone"
				_, err := UnauthenticatedRequester{}.PostForm(client, config, "/oauth/token", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).To(BeNil())
			})
		})

		Describe("PostJson", func() {
			It("calls an endpoint with correct body and headers", func() {
				responseJson = `{ "status" : "great successs" }`

				server.RouteToHandler("POST", "/foo", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("POST", "/foo", ""),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.VerifyJSON(`{"Field1": "hello", "Field2": "world"}`),
				))

				bodyObj := TestData{Field1: "hello", Field2: "world"}

				returnedBytes, _ := UnauthenticatedRequester{}.PostJson(client, config, "/foo", "", bodyObj)
				parsedResponse := string(returnedBytes)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(parsedResponse).To(ContainSubstring("great success"))
			})

			It("returns an error when request fails", func() {
				server.RouteToHandler("POST", "/foo", ghttp.CombineHandlers(
					ghttp.RespondWith(500, "garbage"),
					ghttp.VerifyRequest("POST", "/foo", ""),
				))

				bodyObj := TestData{Field1: "hello", Field2: "world"}
				_, err := UnauthenticatedRequester{}.PostJson(client, config, "/foo", "", bodyObj)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})

			It("supports zone switching", func() {
				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(201, responseJson),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
					ghttp.VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))

				config.ZoneSubdomain = "twilight-zone"
				_, err := UnauthenticatedRequester{}.PostJson(client, config, "/oauth/token", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).To(BeNil())
			})
		})

		Describe("PutJson", func() {
			It("calls an endpoint with correct body and headers", func() {
				responseJson = `{ "status" : "great successs" }`

				server.RouteToHandler("PUT", "/foo", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("PUT", "/foo", ""),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.VerifyJSON(`{"Field1": "hello", "Field2": "world"}`),
				))

				bodyObj := TestData{Field1: "hello", Field2: "world"}

				returnedBytes, _ := UnauthenticatedRequester{}.PutJson(client, config, "/foo", "", bodyObj)
				parsedResponse := string(returnedBytes)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(parsedResponse).To(ContainSubstring("great success"))
			})

			It("returns an error when request fails", func() {
				server.RouteToHandler("PUT", "/foo", ghttp.CombineHandlers(
					ghttp.RespondWith(500, "garbage"),
					ghttp.VerifyRequest("PUT", "/foo", ""),
				))

				bodyObj := TestData{Field1: "hello", Field2: "world"}
				_, err := UnauthenticatedRequester{}.PutJson(client, config, "/foo", "", bodyObj)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})

			It("supports zone switching", func() {
				responseJson = `{ "status" : "great successs" }`

				server.RouteToHandler("PUT", "/foo", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("PUT", "/foo", ""),
					ghttp.VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))

				config.ZoneSubdomain = "twilight-zone"
				UnauthenticatedRequester{}.PutJson(client, config, "/foo", "", TestData{Field1: "hello", Field2: "world"})
				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})
		})
	})

	Describe("AuthenticatedRequester", func() {
		Describe("Get", func() {
			It("calls an endpoint with Accept and Authorization headers", func() {
				server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				config.AddContext(NewContextWithToken("access_token"))
				AuthenticatedRequester{}.Get(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("supports zone switching", func() {
				server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))

				config.AddContext(NewContextWithToken("access_token"))
				config.ZoneSubdomain = "twilight-zone"
				AuthenticatedRequester{}.Get(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("returns a helpful error when GET request fails", func() {
				server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(500, ""),
					ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
				))

				config.AddContext(NewContextWithToken("access_token"))
				_, err := AuthenticatedRequester{}.Get(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})

			It("returns a helpful error when no token in context", func() {
				config.AddContext(NewContextWithToken(""))
				_, err := AuthenticatedRequester{}.Get(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(0))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An access token is required to call"))
			})
		})

		Describe("Delete", func() {
			It("calls an endpoint with Accept and Authorization headers", func() {
				server.RouteToHandler("DELETE", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("DELETE", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				config.AddContext(NewContextWithToken("access_token"))
				AuthenticatedRequester{}.Delete(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("supports zone switching", func() {
				server.RouteToHandler("DELETE", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("DELETE", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))

				config.AddContext(NewContextWithToken("access_token"))
				config.ZoneSubdomain = "twilight-zone"
				AuthenticatedRequester{}.Delete(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("returns a helpful error when DELETE request fails", func() {
				server.RouteToHandler("DELETE", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(500, ""),
					ghttp.VerifyRequest("DELETE", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
				))

				config.AddContext(NewContextWithToken("access_token"))
				_, err := AuthenticatedRequester{}.Delete(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})

			It("returns a helpful error when no token in context", func() {
				config.AddContext(NewContextWithToken(""))
				_, err := AuthenticatedRequester{}.Delete(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(0))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An access token is required to call"))
			})
		})

		Describe("PostForm", func() {
			It("calls an endpoint with correct body and headers", func() {
				responseJson = `{
				  "access_token" : "bc4885d950854fed9a938e96b13ca519",
				  "token_type" : "bearer",
				  "expires_in" : 43199,
				  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
				  "jti" : "bc4885d950854fed9a938e96b13ca519"
				}`

				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
					ghttp.VerifyBody([]byte("hello=world")),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				))

				body := map[string]string{"hello": "world"}
				config.AddContext(NewContextWithToken("access_token"))

				returnedBytes, _ := AuthenticatedRequester{}.PostForm(client, config, "/oauth/token", "", body)
				parsedResponse := string(returnedBytes)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(parsedResponse).To(ContainSubstring("expires_in"))
			})

			It("supports zone switching", func() {
				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
					ghttp.VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))

				config.AddContext(NewContextWithToken("access_token"))
				config.ZoneSubdomain = "twilight-zone"

				AuthenticatedRequester{}.PostForm(client, config, "/oauth/token", "", map[string]string{})
				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("returns an error when request fails", func() {
				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(500, "garbage"),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
				))

				config.AddContext(NewContextWithToken("access_token"))
				_, err := AuthenticatedRequester{}.PostForm(client, config, "/oauth/token", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})

			It("returns a helpful error when no token in context", func() {
				config.AddContext(NewContextWithToken(""))
				_, err := AuthenticatedRequester{}.PostForm(client, config, "/oauth/token", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(0))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An access token is required to call"))
			})
		})

		Describe("PostJson", func() {
			It("calls an endpoint with correct body and headers", func() {
				responseJson = `{ "status" : "great successs" }`

				server.RouteToHandler("POST", "/foo", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("POST", "/foo", ""),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.VerifyJSON(`{"Field1": "hello", "Field2": "world"}`),
				))

				bodyObj := TestData{Field1: "hello", Field2: "world"}
				config.AddContext(NewContextWithToken("access_token"))

				returnedBytes, _ := AuthenticatedRequester{}.PostJson(client, config, "/foo", "", bodyObj)
				parsedResponse := string(returnedBytes)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(parsedResponse).To(ContainSubstring("great success"))
			})

			It("returns an error when request fails", func() {
				server.RouteToHandler("POST", "/foo", ghttp.CombineHandlers(
					ghttp.RespondWith(500, "garbage"),
					ghttp.VerifyRequest("POST", "/foo", ""),
				))

				config.AddContext(NewContextWithToken("access_token"))
				bodyObj := TestData{Field1: "hello", Field2: "world"}
				_, err := AuthenticatedRequester{}.PostJson(client, config, "/foo", "", bodyObj)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})

			It("returns a helpful error when no token in context", func() {
				config.AddContext(NewContextWithToken(""))
				_, err := AuthenticatedRequester{}.PostJson(client, config, "/foo", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(0))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An access token is required to call"))
			})
		})

		Describe("PutJson", func() {
			It("calls an endpoint with correct body and headers", func() {
				responseJson = `{ "status" : "great successs" }`

				server.RouteToHandler("PUT", "/foo", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("PUT", "/foo", ""),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.VerifyJSON(`{"Field1": "hello", "Field2": "world"}`),
				))

				bodyObj := TestData{Field1: "hello", Field2: "world"}
				config.AddContext(NewContextWithToken("access_token"))

				returnedBytes, _ := AuthenticatedRequester{}.PutJson(client, config, "/foo", "", bodyObj)
				parsedResponse := string(returnedBytes)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(parsedResponse).To(ContainSubstring("great success"))
			})

			It("returns an error when request fails", func() {
				server.RouteToHandler("PUT", "/foo", ghttp.CombineHandlers(
					ghttp.RespondWith(500, "garbage"),
					ghttp.VerifyRequest("PUT", "/foo", ""),
				))

				config.AddContext(NewContextWithToken("access_token"))
				bodyObj := TestData{Field1: "hello", Field2: "world"}
				_, err := AuthenticatedRequester{}.PutJson(client, config, "/foo", "", bodyObj)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})

			It("supports zone switching", func() {
				server.RouteToHandler("PUT", "/foo", ghttp.CombineHandlers(
					ghttp.RespondWith(200, `{ "status" : "great successs" }`),
					ghttp.VerifyRequest("PUT", "/foo", ""),
					ghttp.VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))

				config.AddContext(NewContextWithToken("access_token"))
				config.ZoneSubdomain = "twilight-zone"
				_, err := AuthenticatedRequester{}.PutJson(client, config, "/foo", "", TestData{Field1: "hello", Field2: "world"})
				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("returns a helpful error when no token in context", func() {
				config.AddContext(NewContextWithToken(""))
				_, err := AuthenticatedRequester{}.PutJson(client, config, "/foo", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(0))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An access token is required to call"))
			})
		})

	})
})
