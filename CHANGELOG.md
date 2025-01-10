## v1.5.3 (2025-01-10)

### New feature:

- add id to section (#51)([`5a48bcb`](https://github.com/waynezhang/foto/commit/5a48bcb379e9112d33ab73a56c53cb5f41de6901))


## v1.5.2 (2024-11-25)

### New feature:

- #43 Add PNG support([`eb850eb`](https://github.com/waynezhang/foto/commit/eb850eb03967107a73c1f2fc502f57cc793f3bb6))

### Bugs fixed:

- **logger**: init logger in the very beginning([`43862f5`](https://github.com/waynezhang/foto/commit/43862f5d853fa8794747674d1adc0ffd3e86d9f0))


## v1.5.1 (2024-09-27)

### New feature:

- consolidate navigation arrows of lightbox to single option([`003a432`](https://github.com/waynezhang/foto/commit/003a432ff22aab0af6c8c41c57edc41870f215e9))


## v1.5.0 (2024-09-26)

### New feature:

- Add options to customize lightbox([`de1e6e4`](https://github.com/waynezhang/foto/commit/de1e6e4ad1414621e1a404104a39ee78b1a0ddb4))


## v1.4.0 (2024-07-29)

### New feature:

- #34 add `compressQuality` option to configuration file([`c55cfe0`](https://github.com/waynezhang/foto/commit/c55cfe06ae5e1f23e1f6054ac6a42708d23b1952))
- cache versioning([`ea1582a`](https://github.com/waynezhang/foto/commit/ea1582a7e3213f81157d6258468e132581fffcf2))

### Bugs fixed:

- #33 don't show empty section([`62ff944`](https://github.com/waynezhang/foto/commit/62ff944469dab2ddb84854f43e283db57af2d0e0))


## v1.3.0 (2024-06-14)

### New feature:

- support webp format (#25)([`162b155`](https://github.com/waynezhang/foto/commit/162b1559f4de4d8dbe17d941261c8c8136f4f9d9))

## v1.2.2 (2024-02-28)

### Bugs fixed:

- **indexer**: crash on directory traversal([`45de818`](https://github.com/waynezhang/foto/commit/45de8189b028ce71407c422101591cf40835fcb7))

## v1.2.1 (2024-02-07)

### Bugs fixed:

fix(#8): provide iamge size for grid computing ([`35197f6`](https://github.com/waynezhang/foto/commit/35197f60525c52c6f26f5ea460f3ea957e3e6073))
fix: handle singleton in cache properly ([`39c3257`](https://github.com/waynezhang/foto/pull/10/commits/39c32576d4d040647dfa1dbe24ddaffcdd1d8853))
fix: handle relative folders ([`3dc7b96`](https://github.com/waynezhang/foto/commit/3dc7b96d47327c937f695fee5fe911f4bf9b77d1))
fix: show error if failed to clean output directory ([`ec097bf `](https://github.com/waynezhang/foto/commit/ec097bf8002ba27eabbd2a1d2bee372e78e0a11b))

### Others:

refactor: migrate to new logger ([`c075a6a `](https://github.com/waynezhang/foto/commit/c075a6a843045f143bd9202b3cc5ebc16cc3c2fe))

## v1.2.0 (2024-01-05)

### Bugs fixed:

- recursively process directory([`a612b8b`](https://github.com/waynezhang/foto/commit/a612b8bcd9dc4aa825aaef91c54af2c74f4a9264))
- **indexer**: add validation to slug([`0c9fb1e`](https://github.com/waynezhang/foto/commit/0c9fb1ec7a2984db5735905cdbba5700228af2ac))
- **log**: fix some log issue([`29fba6e`](https://github.com/waynezhang/foto/commit/29fba6ee5b326d54b37ed4464b95687bc47b51f2))

### Performance improves:

- build index in concurrency([`7a2b124`](https://github.com/waynezhang/foto/commit/7a2b12417c548e7d79c4c7e1de927369d27def11))
- process images concurrently([`f23b327`](https://github.com/waynezhang/foto/commit/f23b3276c0c59550adc878cab697805939802321))
- none minimizer should do nothing([`07a2b89`](https://github.com/waynezhang/foto/commit/07a2b896f55fd489421e31fc171c96628727e566))

## v1.0.9 (2024-01-02)

### Bugs fixed:

- **image**: size information in generated HTML is incorrect for images([`69072f0`](https://github.com/waynezhang/foto/commit/69072f087cc76adcc4293e473523f869149e4c5c)), Closes: #3

## v1.0.8 (2024-01-01)

### Bugs fixed:

- **export**: add error handling for duplicated slugs([`f378beb`](https://github.com/waynezhang/foto/commit/f378bebacebe96e2ecc31e912492f89200b461ce)), Closes: #4
- HTML should be supported([`7a84081`](https://github.com/waynezhang/foto/commit/7a84081864679040fbda84539c392acae5ad9ec5))

### BREAKING CHANGES:

- PhotoSwipeVersion should be changed to lowercase `photoswipeversion` in the template html.

## v1.0.7 (2023-12-30)

## v1.0.6 (2023-12-30)

## v1.0.5 (2023-08-31)

## v1.0.4 (2022-07-29)
