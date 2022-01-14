package create

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var resposeBody = `
{
    "id": 82706,
    "name": "{name}",
    "type": "{type}",
    "content": "{content}",
    "content_type": "{content_type}"
}
`

func buildResponseContent(req *http.Request) string {
	request := &sdk.CreateResourceRequest{}
	body, _ := ioutil.ReadAll(req.Body)
	_ = json.Unmarshal(body, request)

	response := strings.ReplaceAll(resposeBody, "{name}", request.Name)
	response = strings.ReplaceAll(response, "{type}", request.Trigger)
	response = strings.ReplaceAll(response, "{content}", request.Content)
	response = strings.ReplaceAll(response, "{content_type}", request.ContentType)

	return response
}

func TestCreate(t *testing.T) {
	t.Run("create text resource", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/1234/resources"),
			func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusCreated,
					Request: req,
					Body:    ioutil.NopCloser(strings.NewReader(buildResponseContent(req))),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
		)

		f, stdout, _ := testutils.NewFactory(mock)

		contentFile, _ := os.CreateTemp("", "content.txt")

		_, _ = contentFile.Write([]byte("insert your text here"))

		cmd := NewCmd(f)
		cmd.PersistentFlags().BoolP("verbose", "v", false, "")
		cmd.SetArgs([]string{"-v", "1234", "--name", "/tmp/testando.txt", "--content-type", "text", "--content-file", contentFile.Name()})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, "ID: 82706\nName: /tmp/testando.txt\nType: \nContent type: Text\nContent: \ninsert your text here", stdout.String())

	})

	t.Run("create text resource being verbose", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/1234/resources"),
			func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusCreated,
					Request: req,
					Body:    ioutil.NopCloser(strings.NewReader(buildResponseContent(req))),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
		)

		f, _, _ := testutils.NewFactory(mock)

		contentFile, _ := os.CreateTemp("", "content.txt")

		_, _ = contentFile.Write([]byte("insert your text here"))

		cmd := NewCmd(f)
		cmd.PersistentFlags().BoolP("verbose", "v", false, "")
		cmd.SetArgs([]string{"1234", "--name", "/tmp/testando.txt", "--content-type", "text", "--content-file", contentFile.Name()})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("create script resource", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/1234/resources"),
			func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusCreated,
					Request:    req,
					Body:       ioutil.NopCloser(strings.NewReader(buildResponseContent(req))),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
		)

		f, _, _ := testutils.NewFactory(mock)

		contentFile, _ := os.CreateTemp("", "content.txt")

		_, _ = contentFile.Write([]byte("#!/bin/sh"))

		cmd := NewCmd(f)
		cmd.PersistentFlags().BoolP("verbose", "v", false, "")
		cmd.SetArgs([]string{"1234", "--name", "/tmp/bomb.sh", "--trigger", "Install", "--content-type", "shellscript", "--content-file", contentFile.Name()})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("create script resource without trigger", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		contentFile, _ := os.CreateTemp("", "content.txt")

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234", "--name", "/tmp/bomb.sh", "--content-type", "shellscript", "--content-file", contentFile.Name()})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, utils.ErrorInvalidResourceTrigger)
	})

	t.Run("create resource without content file", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234", "--name", "/tmp/a.txt", "--content-type", "text"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.EqualError(t, err, "required flag(s) \"content-file\" not set")
	})

	t.Run("service not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/1234/resources"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not found"),
		)
		f, _, _ := testutils.NewFactory(mock)

		contentFile, _ := os.CreateTemp("", "content.txt")

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234", "--name", "/tmp/a.txt", "--content-type", "text", "--content-file", contentFile.Name()})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)
		cmd.PersistentFlags().BoolP("verbose", "v", false, "")
		_, err := cmd.ExecuteC()
		require.EqualError(t, err, "Not found. Use -h or --help for more information")
	})
}
