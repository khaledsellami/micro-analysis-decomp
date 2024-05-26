# frozen_string_literal: true

require 'fileutils'
require 'json'

require_relative 'app/utils/cli'
require_relative 'app/analyze'


def main(args)
  app_path = File.absolute_path(args.path)
  output_path = args.output
  logging_level = args.logging
  is_monolithic = args.monolithic
  app_name = File.basename(app_path)
  analyze_app(app_name, app_path, output_path: output_path, logging_level: logging_level, is_monolithic: is_monolithic)
end


if __FILE__ == $PROGRAM_NAME
  # Parse the command line arguments
  args = Cli.parse(ARGV)
  main(args)
end
