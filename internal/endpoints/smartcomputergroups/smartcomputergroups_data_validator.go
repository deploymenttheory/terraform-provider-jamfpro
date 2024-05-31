// smartcomputergroup_data_validator.go
package smartcomputergroups

// getCriteriaOperators returns a list of criteria operators for Smart Computer Groups.
func getCriteriaOperators() []string {
	var out []string
	out = []string{
		And,
		Or,
		SearchTypeIs,
		SearchTypeIsNot,
		SearchTypeHas,
		SearchTypeDoesNotHave,
		SearchTypeMemberOf,
		SearchTypeNotMemberOf,
		SearchTypeBeforeYYYYMMDD,
		SearchTypeAfterYYYYMMDD,
		SearchTypeMoreThanXDaysAgo,
		SearchTypeLessThanXDaysAgo,
		SearchTypeLike,
		SearchTypeNotLike,
		SearchTypeGreaterThan,
		SearchTypeMoreThan,
		SearchTypeLessThan,
		SearchTypeGreaterThanOrEqual,
		SearchTypeLessThanOrEqual,
		SearchTypeMatchesRegex,
		SearchTypeDoesNotMatch,
	}
	return out
}
