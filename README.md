# go.bug.st/relaxed-semver [![build status](https://github.com/bugst/relaxed-semver/workflows/test/badge.svg)](https://travis-ci.org/bugst/relaxed-semver) [![codecov](https://codecov.io/gh/bugst/relaxed-semver/branch/master/graph/badge.svg)](https://codecov.io/gh/bugst/relaxed-semver)

A library for handling a superset of semantic versioning in golang.

## Documentation and examples

See the godoc here: https://godoc.org/go.bug.st/relaxed-semver

## Semantic versioning specification followed in this library

This library tries to implement the semantic versioning specification [2.0.0](https://semver.org/spec/v2.0.0.html) with an exception: the numeric format `major.minor.patch` like `1.3.2` may be truncated if a number is zero, so:

- `1.2.0` or `1.2.0-beta` may be written as `1.2` or `1.2-beta` respectively
- `1.0.0` or `1.0.0-beta` may be written `1` or `1-beta` respectively
- `0.0.0` may be written as the **empty string**, but `0.0.0-beta` may **not** be written as `-beta`

## Usage

You can parse a semver version string with the `Parse` function that returns a `Version` object that can be used to be compared with other `Version` objects using the `CompareTo`, `LessThan` , `LessThanOrEqual`, `Equal`, `GreaterThan` and `GreaterThanOrEqual` methods.

The `Parse` function returns an `error` if the string does not comply to the above specification. Alternatively the `MustParse` function can be used, it returns only the `Version` object or panics if a parsing error occurs.

## Why Relaxed?

This library allows the use of an even more relaxed semver specification using the `RelaxedVersion` object. It works with the following rules:

- If the parsed string is a valid semver (following the rules above), then the `RelaxedVersion` will behave exactly as a normal `Version` object
- if the parsed string is **not** a valid semver, then the string is kept as-is inside the `RelaxedVersion` object as a custom version string
- when comparing two `RelaxedVersion` the rule is simple: if both are valid semver, the semver rules applies; if both are custom version string they are compared as alphanumeric strings; if one is valid semver and the other is a custom version string the valid semver is always greater
- two `RelaxedVersion` are compatible (by the `CompatibleWith` operation) only if
  - they are equal
  - they are both valid semver and they are compatible as per semver specification

The `RelaxedVersion` object is basically made to allow systems that do not use semver to soft transition to semantic versioning, because it allows an intermediate period where the invalid version is still tolerated.

To parse a `RelaxedVersion` you can use the `ParseRelaxed` function.

## Version constraints

Dependency version matching can be specified via version constraints, which might be a version range or an exact version.

The following operators are supported:

|          |                          |
| -------- | ------------------------ |
| `=`      | equal to                 |
| `>`      | greater than             |
| `>=`     | greater than or equal to |
| `<`      | less than                |
| `<=`     | less than or equal to    |
| `^`      | compatible-with          |
| `!`      | NOT                      |
| `&&`     | AND                      |
| `\|\|`   | OR                       |
| `(`, `)` | constraint group         |

### Examples

Given the following releases of a dependency:

- `0.1.0`
- `0.1.1`
- `0.2.0`
- `1.0.0`
- `2.0.0`
- `2.0.5`
- `2.0.6`
- `2.1.0`
- `3.0.0`

constraints conditions would match as follows:

| The following condition...       | will match with versions...                                            |
| -------------------------------- | ---------------------------------------------------------------------- |
| `=1.0.0`                         | `1.0.0`                                                                |
| `>1.0.0`                         | `2.0.0`, `2.0.5`, `2.0.6`, `2.1.0`, `3.0.0`                            |
| `>=1.0.0`                        | `1.0.0`, `2.0.0`, `2.0.5`, `2.0.6`, `2.1.0`, `3.0.0`                   |
| `<2.0.0`                         | `0.1.0`, `0.1.1`, `0.2.0`, `1.0.0`                                     |
| `<=2.0.0`                        | `0.1.0`, `0.1.1`, `0.2.0`, `1.0.0`, `2.0.0`                            |
| `!=1.0.0`                        | `0.1.0`, `0.1.1`, `0.2.0`, `2.0.0`, `2.0.5`, `2.0.6`, `2.1.0`, `3.0.0` |
| `>1.0.0 && <2.1.0`               | `2.0.0`, `2.0.5`, `2.0.6`                                              |
| `<1.0.0 \|\| >2.0.0`             | `0.1.0`, `0.1.1`, `0.2.0`, `2.0.5`, `2.0.6`, `2.1.0`, `3.0.0`          |
| `(>0.1.0 && <2.0.0) \|\| >2.0.5` | `0.1.1`, `0.2.0`, `1.0.0`, `2.0.6`, `2.1.0`, `3.0.0`                   |
| `^2.0.5`                         | `2.0.5`, `2.0.6`, `2.1.0`                                              |
| `^0.1.0`                         | `0.1.0`, `0.1.1`                                                       |

## Json parsable

The `Version` and `RelaxedVersion` have the JSON un/marshaler implemented so they can be JSON decoded/encoded.

## Binary/GOB encoding support

The `Version` and `RelaxedVersion` provides optimized `MarshalBinary`/`UnmarshalBinary` methods for binary encoding.

## Yaml parsable with `gopkg.in/yaml.v3`

The `Version` and `RelaxedVersion` have the YAML un/marshaler implemented so they can be YAML decoded/encoded with the excellent `gopkg.in/yaml.v3` library.

## SQL support

The `Version` and `RelaxedVersion` types provides the `sql.Scanner` and `driver.Valuer` interfaces. Those objects could be directly used in SQL queries, their value will be mapped into a string field.

## Lexicographic sortable strings that keeps semantic versioning order

The `Version` and `RelaxedVersion` objects provides the `SortableString()` method that returns a string with a peculiar property: the alphanumeric sorting of two `Version.SortableString()` matches the semantic versioning ordering of the underling `Version` objects. In other words, given two `Version` object `a` and `b`:
* if `a.LessThan(b)` is true then `a.SortableString() < b.SortableString()` is true and vice-versa.
* if `a.Equals(b)` is true then `a.SortableString() == b.SortableString()` is true and vice-versa.
* more generally, the following assertion is always true: `a.CompareTo(b) == cmp.Compare(a.SortableString(), b.SortableString())`

This is accomplished by adding some adjustment characters to the original semver `Version` string with the purpose to change the behaviour of the natural alphabetic ordering, in particular:
* to allow comparison of numeric values (keeping digits aligned by unit, tenths, hundhereds, etc...).
* to invert the ordering of versions with and without prelease (a version with prelease should be lower priority compared to the same version without prerelease, but being longer alphanumerically it naturally follows it).

To give you an idea on how it works, the following table shows some examples of semver versions and their `SortableString` counter part:

| semver             | `SortableString()`  |
| ------------------ | ------------------- |
| `1.2.4`            | `;1.2.4;`           |
| `1.3.0-rc`         | `;1.3.0-;rc`        |
| `1.3.0-rc.0`       | `;1.3.0-;rc,:0`     |
| `1.3.0-rc.5`       | `;1.3.0-;rc,:5`     |
| `1.3.0-rc.5+build` | `;1.3.0-;rc,:5`     |
| `1.3.0-rc.20`      | `;1.3.0-;rc,::20`   |
| `1.3.0-rc-`        | `;1.3.0-;rc-`       |
| `1.3.0`            | `;1.3.0;`           |
| `1.20.0`           | `;1.:20.0;`         |
| `1.90.0`           | `;1.:90.0;`         |
| `1.300.0-6`        | `;1.::300.0-:6`     |
| `1.300.0-30`       | `;1.::300.0-::30`   |
| `1.300.0-1pre`     | `;1.::300.0-;1pre`  |
| `1.300.0-pre`      | `;1.::300.0-;pre`   |
| `1.300.0`          | `;1.::300.0;`       |

The `SortableString()` can be used in SQL databases to simplify the ordering of a set of versions in a table.
