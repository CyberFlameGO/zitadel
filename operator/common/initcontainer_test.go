package common

import (
	"github.com/caos/zitadel/operator/helpers"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestDeployment_GetInitContainer(t *testing.T) {
	users := []string{"test"}
	version := "test"

	initCommands := []string{
		"cp /dbsecrets/client_root/ca.crt /tmp/dbsecrets/ca.crt",
		"cp /dbsecrets/client_test/client.test.crt /tmp/dbsecrets/client.test.crt",
		"cp /dbsecrets/client_test/client.test.key /tmp/dbsecrets/client.test.key",
		//"chown -R 1000:1000 /tmp/dbsecrets",
		"chmod 0600 /tmp/dbsecrets/*",
	}

	initVolumeMounts := []corev1.VolumeMount{
		{Name: rootSecret, MountPath: certMountPath + "/client_root"},
		{Name: dbSecrets, MountPath: certTempMountPath},
		{Name: "client-test", MountPath: certMountPath + "/client_test"},
	}

	equals := corev1.Container{
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:  helpers.PointerInt64(1000),
			RunAsGroup: helpers.PointerInt64(1000),
		},
		Name:                     "fix-permissions-zitadel",
		Image:                    ZITADELCockroachImage.Reference("", version),
		Command:                  []string{"/bin/sh", "-c"},
		Args:                     []string{strings.Join(initCommands, " && ")},
		VolumeMounts:             initVolumeMounts,
		ImagePullPolicy:          corev1.PullIfNotPresent,
		TerminationMessagePolicy: "File",
		TerminationMessagePath:   "/dev/termination-log",
	}

	init := GetInitContainer("zitadel", rootSecret, dbSecrets, users, ZITADELCockroachImage.RunAsUser(), ZITADELCockroachImage.Reference("", version))

	assert.Equal(t, equals, init)
}

func TestDeployment_GetInitContainer1(t *testing.T) {
	users := []string{"test1"}
	version := "test1"

	initCommands := []string{
		"cp /dbsecrets/client_root/ca.crt /tmp/dbsecrets/ca.crt",
		"cp /dbsecrets/client_test1/client.test1.crt /tmp/dbsecrets/client.test1.crt",
		"cp /dbsecrets/client_test1/client.test1.key /tmp/dbsecrets/client.test1.key",
		//"chown -R 1000:1000 /tmp/dbsecrets",
		"chmod 0600 /tmp/dbsecrets/*",
	}

	initVolumeMounts := []corev1.VolumeMount{
		{Name: rootSecret, MountPath: certMountPath + "/client_root"},
		{Name: dbSecrets, MountPath: certTempMountPath},
		{Name: "client-test1", MountPath: certMountPath + "/client_test1"},
	}

	equals := corev1.Container{
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:  helpers.PointerInt64(1000),
			RunAsGroup: helpers.PointerInt64(1000),
		},
		Name:                     "fix-permissions-zitadel",
		Image:                    ZITADELCockroachImage.Reference("", version),
		Command:                  []string{"/bin/sh", "-c"},
		Args:                     []string{strings.Join(initCommands, " && ")},
		VolumeMounts:             initVolumeMounts,
		TerminationMessagePolicy: "File",
		TerminationMessagePath:   "/dev/termination-log",
		ImagePullPolicy:          corev1.PullIfNotPresent,
	}

	init := GetInitContainer("zitadel", rootSecret, dbSecrets, users, ZITADELCockroachImage.RunAsUser(), ZITADELCockroachImage.Reference("", version))

	assert.Equal(t, equals, init)
}

func TestDeployment_GetInitContainer2(t *testing.T) {
	users := []string{"test1", "test2"}
	version := "test2"

	initCommands := []string{
		"cp /dbsecrets/client_root/ca.crt /tmp/dbsecrets/ca.crt",
		"cp /dbsecrets/client_test1/client.test1.crt /tmp/dbsecrets/client.test1.crt",
		"cp /dbsecrets/client_test1/client.test1.key /tmp/dbsecrets/client.test1.key",
		"cp /dbsecrets/client_test2/client.test2.crt /tmp/dbsecrets/client.test2.crt",
		"cp /dbsecrets/client_test2/client.test2.key /tmp/dbsecrets/client.test2.key",
		//"chown -R 1000:1000 /tmp/dbsecrets",
		"chmod 0600 /tmp/dbsecrets/*",
	}

	initVolumeMounts := []corev1.VolumeMount{
		{Name: rootSecret, MountPath: certMountPath + "/client_root"},
		{Name: dbSecrets, MountPath: certTempMountPath},
		{Name: "client-test1", MountPath: certMountPath + "/client_test1"},
		{Name: "client-test2", MountPath: certMountPath + "/client_test2"},
	}

	equals := corev1.Container{
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:  helpers.PointerInt64(1000),
			RunAsGroup: helpers.PointerInt64(1000),
		},
		Name:                     "fix-permissions-zitadel",
		Image:                    ZITADELCockroachImage.Reference("", version),
		Command:                  []string{"/bin/sh", "-c"},
		Args:                     []string{strings.Join(initCommands, " && ")},
		VolumeMounts:             initVolumeMounts,
		ImagePullPolicy:          corev1.PullIfNotPresent,
		TerminationMessagePolicy: "File",
		TerminationMessagePath:   "/dev/termination-log",
	}

	init := GetInitContainer("zitadel", rootSecret, dbSecrets, users, ZITADELCockroachImage.RunAsUser(), ZITADELCockroachImage.Reference("", version))

	assert.Equal(t, equals, init)
}