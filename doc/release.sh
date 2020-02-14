#!/bin/bash - 
#===============================================================================
#
#          FILE: release.sh
#
#         USAGE: ./release.sh <VERSION>
#
#   DESCRIPTION: This script re-generates the main index.html (=> README.md)
#                with a reference to the latest version URL.
#
#       OPTIONS: ---
#  REQUIREMENTS: cat
#          BUGS: ---
#         NOTES: ---
#        AUTHOR: Francesco Emanuel Bennici <benniciemanuel78@gmail.com>
#  ORGANIZATION: FABMation GmbH
#       CREATED: 02/14/2020 05:45:19 PM
#      REVISION: 001
#===============================================================================

set -o nounset                              # Treat unset variables as an error

# get directory of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "---
meta:
  - http-equiv: refresh
    content: 0;url=/${1}
---" > "${DIR}/README.md"
