# pass-man

Terminal-based multi-vault password manager application built for the purpose of understanding Go better.
This application makes use of *AES-256 GCM* encryption and *Argon2* key derivation function for secure local storage of vault credentials, supported through Go's `crypto` module.

Inspired by [Coding Challenges #58](https://codingchallenges.substack.com/p/coding-challenge-58-password-manager)

## Setup Instructions

1. Ensure that the latest version of Go (>=1.22.1) is installed. Go can be installed from [here](https://go.dev/dl/).
2. Clone this repo using the following command:
```
git clone https://github.com/alpha2303/pass-man.git
```
2. Once cloned, open a terminal at the root folder of this repo and run `go mod download` to download the necessary dependencies.
3. To verify that the dependencies have been downloaded without issues, run `go mod verify`.

## Run Application in "Dev" mode

To run the application without explicitly building it (similar to dev mode), run `go run .` in the terminal opened at the root folder of the repo.

## Build Instructions
1. Ensure that all the necessary dependency modules have been downloaded without issues.
2. In the terminal opened at the root folder, run `go build` to build the application. The application will appear in the same folder as an executable file named `pass-man`.
