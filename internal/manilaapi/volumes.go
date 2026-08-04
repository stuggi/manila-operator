package manilaapi

import (
	manilav1 "github.com/openstack-k8s-operators/manila-operator/api/v1beta1"
	"github.com/openstack-k8s-operators/manila-operator/internal/manila"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumes -
func GetVolumes(parentName string, name string, extraVol []manilav1.ManilaExtraVolMounts) []corev1.Volume {
	var config0644AccessMode int32 = 0644

	apiVolumes := []corev1.Volume{
		{
			Name: "config-data-custom",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: &config0644AccessMode,
					SecretName:  name + "-config-data",
				},
			},
		},
		GetRunHttpdVolume(),
		GetVarLogHttpdVolume(),
	}

	return append(manila.GetVolumes(parentName, extraVol, manila.ManilaAPIPropagation), apiVolumes...)
}

// GetVolumeMounts - ManilaAPI VolumeMounts
func GetVolumeMounts(extraVol []manilav1.ManilaExtraVolMounts) []corev1.VolumeMount {
	apiVolumeMounts := []corev1.VolumeMount{
		{
			Name:      "config-data-custom",
			MountPath: "/etc/manila/manila.conf.d",
			ReadOnly:  true,
		},
		{
			Name:      "config-data",
			MountPath: "/etc/httpd/conf/httpd.conf",
			SubPath:   "httpd.conf",
			ReadOnly:  true,
		},
		{
			Name:      "config-data",
			MountPath: "/etc/httpd/conf.d/10-manila_wsgi.conf",
			SubPath:   "10-manila_wsgi.conf",
			ReadOnly:  true,
		},
		GetRunHttpdVolumeMount(),
		GetVarLogHttpdVolumeMount(),
	}

	return append(manila.GetVolumeMounts(extraVol, manila.ManilaAPIPropagation), apiVolumeMounts...)
}

// GetRunHttpdVolume - EmptyDir for httpd's PID file directory, needed once
// httpd runs as a non-root, FSGroup-only user (kolla used to chown
// /etc/httpd/run at startup).
func GetRunHttpdVolume() corev1.Volume {
	return corev1.Volume{
		Name: "run-httpd",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
}

// GetRunHttpdVolumeMount -
func GetRunHttpdVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      "run-httpd",
		MountPath: "/etc/httpd/run",
	}
}

// GetVarLogHttpdVolume - EmptyDir for httpd's own logs, defense-in-depth for
// any RPM-shipped conf.d file that references a relative "logs/*" path (this
// service's own ErrorLog/CustomLog already go to /dev/stdout).
func GetVarLogHttpdVolume() corev1.Volume {
	return corev1.Volume{
		Name: "var-log-httpd",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
}

// GetVarLogHttpdVolumeMount -
func GetVarLogHttpdVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      "var-log-httpd",
		MountPath: "/var/log/httpd",
	}
}

// GetLogVolumeMount - Manila API LogVolumeMount
func GetLogVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      logVolume,
		MountPath: "/var/log/manila",
		ReadOnly:  false,
	}
}

// GetLogVolume - Manila API LogVolume
func GetLogVolume() corev1.Volume {
	return corev1.Volume{
		Name: logVolume,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{Medium: ""},
		},
	}
}
