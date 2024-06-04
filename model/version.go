package model

import (
	"errors"
	"strconv"
	"strings"
)

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

func (v *Version) SetMajor(major int) {
	v.Major = major
}

func (v *Version) SetMinor(minor int) {
	v.Minor = minor
}

func (v *Version) SetPatch(patch int) {
	v.Patch = patch
}

func ToString(ver Version) string {
	return strconv.Itoa(ver.Major) + "." + strconv.Itoa(ver.Minor) + "." + strconv.Itoa(ver.Patch)
}

func ToVersion(version string) (*Version, error) {
	split := strings.Split(version, ".")
	if len(split) != 3 {
		return &Version{}, errors.New("version incorrect")
	}

	major, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, errors.New("failed converting")
	}
	minor, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, errors.New("failed converting")
	}
	patch, err := strconv.Atoi(split[2])
	if err != nil {
		return nil, errors.New("failed converting")
	}

	versionModel := Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}

	return &versionModel, nil
}

type VersionRepository interface {
	Delete()
	Update()
	Create()
	Find()
}
