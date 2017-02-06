package commands

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/premkit/premkit/server"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCkJf3/LhtnSy1p
2iTEAnn93eFfm4U8AiLhk33ipQZJHXYO0RcRhw5jQSoY1G2gqFMrbYA69nTOg0ez
OY7L8LI1BUkMw5qA1XychfXIYy8ZvQVObaOy2W6HKe/wd+3LRjbjir2fq9nehdaB
mENTig3MF6ryLaaSnRY4F2Pg1lf9EXjPdeuoQqhg5weFh+xJW/E+zgRoghwz2nya
HEckEAXYU8ZORB7HVh+oVVeggIU14ogawHK4vdRFCHzkn6FisacBDMrlszOASWiS
9ncv0q3p4s27Kb2PIhuOjHQJ+0EdCcceGcqtYlq/61kbkPH7ISpf8bBtK3BnatSI
xS2QJQm1AgMBAAECggEANRQz9fgq1FPy82+ew+MpH3ZIEmpvwt/N97OB2XATgEEO
k+v40aoidOX1fuHyMSk8+6YE+QwI6V56KPJLwpaqiYqT/JSjuVVPXi3TNGEeMex1
cs7xSDwXCY3+EHw3YKvrw9hxSNiBMvuESZO68aCKpZxhor4wRuiU7r5hharJ+QLQ
ti+2EFIfUWwZS77Q5auBGgSKdP+bf3HMB3xsSkH3DMBvIqcQOcKDvGZ6JMYx+KDc
x5OUYxcSnCM8N9C4e4SzkF7Ej791DFszyBnSDiXR1WjJf5skIymgWbManM5O2u82
S0119C9roh9rv8UQc0P4ke6EFq2kOzW/BNRBD1wlYQKBgQDSr5xHa9h2HrepCKcT
u5WdXt4BZVJQbsMzds0ZfWzRCeGzVvE34dXnF1zGhfaN5hmg6BOEpHM4ERJU0SXE
d1ZKyjZ1CWSnjwaukbUHQjiUDiMdOvBddu4V2CuODcYmBM/7RiqaPpbJSmFOsaHh
qTjRVhxoBch6W0pEjikyj8rs1wKBgQDHdATJ1CNESZXGYMEwP7udi2fjZHp+LeG0
fV9eRuNcDtBx1FyguclPrSnGN24mGYkOI3+njEbedZs2uoHhaQXw+poqY/Dw4qo6
A2khXCjvfa+pDgk0DA212O8ZyRLZ2PMjLQduQz3+EGqohuqE3TWdGTXA8riMmkLX
wGRwFe/AUwKBgQC9PYOQG2x43Kp3KBB6hvmiOv4KHupK2NJ4vXMIPEKrmMakAan1
WeJ6CeAJaXbGijHm983gTJ45dAwVJy9XQyG9V9iGU4OXhb6ourPx6ydKxVABB1mz
egnskRi+Jd0fdR8jQikuFp31+9tfheoz+X3RehlVziv+y1TwMwkKI2JQTQKBgHRY
S/bDhTLvTavjgq23b6SNzjMJyJ5T+0YCoB/pb/SiO5s6yjGDTlfo5eZXLSySVq1l
rbA5lplrtveswdiQH8QbGtTBaanKPowKs0efb82L3mzZ4Cp5IYJDIe5DqXhkIigR
uzTpin7qap0V3jVUqFKUgxOjQl3aGkWqV6w+T5U7AoGAIewvru6ygrLU3jyCzBO0
flaSyH71AsVXdqCBghzPsNWHoVDsu3JQ8IcOqc59yjxtcNK9rGFkUPSagP4A6yt5
IGYXk3u+UivLXuDFzEkMA2FBMYbdxGvJMeeTda2aBV/uVD2JjPc1ahamfHpX6QVx
W1xCZvCDq1V1YC0pS4Q82jg=
-----END PRIVATE KEY-----`

	testCert = `-----BEGIN CERTIFICATE-----
MIIDszCCApugAwIBAgIJAJba86Dn9p9CMA0GCSqGSIb3DQEBCwUAMHAxCzAJBgNV
BAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRQwEgYDVQQHDAtMb3MgQW5nZWxl
czEQMA4GA1UECgwHVEVTVElORzEQMA4GA1UECwwHVEVTVElORzESMBAGA1UEAwwJ
bG9jYWxob3N0MB4XDTE2MDcxMzE4NTQ0OVoXDTE3MDcxMzE4NTQ0OVowcDELMAkG
A1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFDASBgNVBAcMC0xvcyBBbmdl
bGVzMRAwDgYDVQQKDAdURVNUSU5HMRAwDgYDVQQLDAdURVNUSU5HMRIwEAYDVQQD
DAlsb2NhbGhvc3QwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCkJf3/
LhtnSy1p2iTEAnn93eFfm4U8AiLhk33ipQZJHXYO0RcRhw5jQSoY1G2gqFMrbYA6
9nTOg0ezOY7L8LI1BUkMw5qA1XychfXIYy8ZvQVObaOy2W6HKe/wd+3LRjbjir2f
q9nehdaBmENTig3MF6ryLaaSnRY4F2Pg1lf9EXjPdeuoQqhg5weFh+xJW/E+zgRo
ghwz2nyaHEckEAXYU8ZORB7HVh+oVVeggIU14ogawHK4vdRFCHzkn6FisacBDMrl
szOASWiS9ncv0q3p4s27Kb2PIhuOjHQJ+0EdCcceGcqtYlq/61kbkPH7ISpf8bBt
K3BnatSIxS2QJQm1AgMBAAGjUDBOMB0GA1UdDgQWBBTopUwpFkAdPQOnTIFjb17N
0vLBpDAfBgNVHSMEGDAWgBTopUwpFkAdPQOnTIFjb17N0vLBpDAMBgNVHRMEBTAD
AQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBhhWok7ZA/49iLlCdrhhIM4avwZPq0y+S5
omt0RKhyxUMFLZaRBlxTmFcfbBlrZs4OSDjPvhOhGwk8yCI9IjllZvhqdY//v7kg
ok+ozFGdGf12J1DqBAyEh4hFI3WpYsgrsvnOdOtfkgpPbGqYbntjyyQa5WIXwZO3
lzyk5IGsswArL97CuNJcGD2EVAuu/nv3gwKLJu5Zk2ed4TPMhC+afMp1MxJZVVrZ
MLCFw1+qcNqqRWj7vtOqL7JTVInIDDdMhQNFL1vZHi3ALJHJKZ9jwoot3XlFJjpF
qO2mGEQNtAsgh5s827CAveJdU+FmdZu2cUIqkCtPGC2fURpGB9vN
-----END CERTIFICATE-----`
)

func TestBuildConfigFromFile(t *testing.T) {
	// Write the test key and cert to temp files
	dirName, err := ioutil.TempDir("", "premkit-test")
	require.NoError(t, err)
	defer os.RemoveAll(dirName)

	err = ioutil.WriteFile(path.Join(dirName, "key"), []byte(testKey), 0400)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(dirName, "cert"), []byte(testCert), 0655)
	require.NoError(t, err)

	// Set up the test
	viper.Set("key_file", path.Join(dirName, "key"))
	viper.Set("cert_file", path.Join(dirName, "cert"))
	viper.Set("self_signed", false)

	config, err := buildConfig()
	require.NoError(t, err)
	assert.NotNil(t, config)
	expectedConfig := server.Config{
		HTTPPort:    2080,
		HTTPSPort:   2443,
		TLSKeyFile:  path.Join(dirName, "key"),
		TLSCertFile: path.Join(dirName, "cert"),
	}
	assert.Equal(t, expectedConfig, *config)
}
