# typed: true
# frozen_string_literal: true

$CACHE_DIR = File.join(Dir.home, '.cache/rpm-get')

# Creates the cache dir.
def create_cache_dir = Dir.mkdir($CACHE_DIR)

# create_cache_dir()
