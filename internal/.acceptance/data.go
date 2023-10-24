// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

const charSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"

type TestData struct {
	RandomInteger int    // Random integer unique to this test case
	RandomString  string // Random 5-character string unique to this test case
	ResourceName  string // Fully qualified resource name
	ResourceType  string // Terraform Resource Type
	resourceLabel string // Label used for the resource, generally "test"
}

// BuildTestData generates test data for the given resource
func BuildTestData(t *testing.T, resourceType string, resourceLabel string) TestData {
	testData := TestData{
		RandomInteger: RandTimeInt(),
		RandomString:  randString(5),
		ResourceName:  fmt.Sprintf("%s.%s", resourceType, resourceLabel),
		ResourceType:  resourceType,
		resourceLabel: resourceLabel,
	}

	return testData
}

// RandTimeInt generates a random integer based on the current time
func RandTimeInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()
}

// RandomIntOfLength returns a random integer of the specified length
func (td *TestData) RandomIntOfLength(len int) int {
	if len < 8 || len > 18 {
		panic("Invalid Test: RandomIntOfLength: len is not between 8 or 18 inclusive")
	}

	if len >= 16 {
		return td.RandomInteger / int(math.Pow10(18-len))
	}

	s := strconv.Itoa(td.RandomInteger)
	r := s[16:18]
	v := s[0 : len-2]
	i, _ := strconv.Atoi(v + r)

	return i
}

// RandomStringOfLength returns a random string of the specified length
func (td *TestData) RandomStringOfLength(len int) string {
	if len < 1 || len > 1024 {
		panic("Invalid Test: RandomStringOfLength: length argument must be between 1 and 1024 characters")
	}

	return randString(len)
}

// randString generates a random alphanumeric string of the length specified
func randString(strlen int) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSetAlphaNum[rand.Intn(len(charSetAlphaNum))]
	}
	return string(result)
}
