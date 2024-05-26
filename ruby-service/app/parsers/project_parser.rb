# frozen_string_literal: true

require_relative 'file_parser'

class ProjectParser
  EXCLUDED = %w[vendor lib]

  attr_reader :objects_, :executables_

  def initialize(root_path, service_name)
    @root_path = root_path
    @service_name = service_name
    @objects_ = []
    @executables_ = []
  end

  def find_files
    files = []
    Dir.glob(File.join(@root_path, '**', '*.rb')) do |file|
      files << file if EXCLUDED.none? { |dir| File.fnmatch(File.join('*', dir, '*'), file) }
    end
    files
  end

  def parse
    ruby_files = find_files
    ruby_files.each do |file|
      analyzer = FileParser.new(@root_path, file, @service_name)
      results = analyzer.start
      @objects_ += results[0]
      @executables_ += results[1]
    end
  end
end
