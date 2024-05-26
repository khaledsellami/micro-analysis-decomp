# frozen_string_literal: true

require 'rubocop-ast'

require_relative 'node_processor'
require_relative '../utils/multi_logger'


class FileParser
  def initialize(root_path, file_path, service_name)
    @file_path = file_path
    @ruby_version = 3.4 # should look into a way to get the ruby version from the file
    @root_path = root_path
    @service_name = service_name
    @logger = MultiLogger.instance
  end

  def start
    code = File.read(@file_path)
    source = RuboCop::AST::ProcessedSource.new(code, @ruby_version)
    processor = NodeProcessor.new
    processor.init_vars(@root_path, @file_path, @service_name)
    processor.process(source.ast)
    [processor.objects_, processor.executables_]
  end
end
