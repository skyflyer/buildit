package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAllYAML(t *testing.T) {
	content := []byte(`repo: ena dva tri
steps:
- prvic
- drugic s presledki
- tretjic
`)
	c, err := parseYaml(content)
	require.Nil(t, err)
	require.Equal(t, "ena dva tri", c.Repo)
	require.Equal(t, 3, len(c.Steps))
	require.Equal(t, "prvic", c.Steps[0])
	require.Equal(t, "drugic s presledki", c.Steps[1])
	require.Equal(t, "tretjic", c.Steps[2])
}

func TestParseEmpty(t *testing.T) {
	_, err := parseYaml([]byte(``))

	require.NotNil(t, err, "Empty config should return error")
}

func TestRequireRepo(t *testing.T) {
	content := []byte(`steps:
- prvic
- drugic s presledki
- tretjic
`)
	c, err := parseYaml(content)
	require.Equal(t, 3, len(c.Steps))
	require.NotNil(t, err, "Repo must be required")
}

func TestRequireSteps(t *testing.T) {
	content := []byte(`repo: ena dva tri`)

	c, e := parseYaml(content)
	require.Equal(t, "ena dva tri", c.Repo)
	require.NotNil(t, e, "Steps are required")
}

func TestDoesNotRequireBranch(t *testing.T) {
	content := []byte(`repo: yes
steps:
- fake
`)

	c, e := parseYaml(content)

	require.Nil(t, e)
	require.Empty(t, c.Branch)
}

func TestBranchParsedCorrectly(t *testing.T) {
	content := []byte(`repo: yes
branch: develop
steps:
- fake
`)

	c, e := parseYaml(content)

	require.Nil(t, e)
	require.Equal(t, "develop", c.Branch)
}

func TestParseNoAuth(t *testing.T) {
	content := []byte(`repo: yes
branch: develop
steps:
- fake
`)
	c, e := parseYaml(content)

	require.Nil(t, e)
	require.Nil(t, c.Auth)
}

func TestParseAuthWithPubkey(t *testing.T) {
	content := []byte(`repo: yes
branch: develop
steps:
- fake
auth:
  key: /Users/miha/.ssh/id_rsa
`)
	c, e := parseYaml(content)

	require.Nil(t, e)
	require.NotNil(t, c.Auth)
	require.Equal(t, "/Users/miha/.ssh/id_rsa", c.Auth.Key)
}

func TestParseAuthUserPwd(t *testing.T) {
	content := []byte(`repo: yes
branch: develop
steps:
- fake
auth:
  username: u
  password:
`)
	c, e := parseYaml(content)

	require.Nil(t, e)
	require.NotNil(t, c.Auth)
	require.Equal(t, "u", c.Auth.Username)
	require.Equal(t, "", c.Auth.Password)
}

func TestParseOptionalWorkDir(t *testing.T) {
	content := []byte(`repo: yes
branch: develop
steps:
- fake
auth:
  username: u
  password:
`)
	c, e := parseYaml(content)

	require.Nil(t, e)
	require.NotNil(t, c.WorkingDirectory)
	require.Equal(t, DefaultWorkingDirectory, c.WorkingDirectory)
}

func TestParseDefinedWorkDir(t *testing.T) {
	content := []byte(`repo: yes
workdir: someplace
branch: develop
steps:
- fake
auth:
  username: u
  password:
`)
	c, e := parseYaml(content)

	require.Nil(t, e)
	require.NotNil(t, c.WorkingDirectory)
	require.Equal(t, "someplace", c.WorkingDirectory)
}

func TestParseAuthNoSSHAgent(t *testing.T) {
	content := []byte(`repo: yes
workdir: someplace
branch: develop
steps:
- fake
auth:
  username: u
  password:
`)
	c, e := parseYaml(content)

	require.Nil(t, e)
	require.Equal(t, false, c.Auth.UseSSHAgent)
}

func TestParseAuthSSHAgent(t *testing.T) {
	content := []byte(`repo: yes
workdir: someplace
branch: develop
steps:
- fake
auth:
  username: u
  password:
  use_ssh_agent: true
`)
	c, e := parseYaml(content)

	require.Nil(t, e)
	require.Equal(t, true, c.Auth.UseSSHAgent)
}
