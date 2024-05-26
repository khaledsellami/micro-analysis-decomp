# frozen_string_literal: true
require 'pathname'

require 'rubocop-ast'

require_relative '../models/object_'
require_relative '../models/executable_'
require_relative '../utils/multi_logger'


class NodeProcessor < Parser::AST::Processor
  include RuboCop::AST::Traversal

  attr_reader :objects_, :executables_

  PACKAGE_SEPARATOR = '/'
  PACKAGE_NAME_SEPARATOR = '::'
  PARENT_SEPARATOR = '.'
  MAX_CLASS_DEPTH = 0
  MAX_METHOD_DEPTH = 1

  def initialize
    @root_path = nil
    @file_path = nil
    @service_name = nil
    @package_simple_name = nil
    @package_name = nil
    @objects_ = []
    @executables_ = []
    @depth = 0
    @method_depth = 0
    @current_parent = nil
    @last_class = nil
    @logger = MultiLogger.instance
    @initialized = false
    super
  end

  def init_vars(root_path, file_path, service_name)
    @root_path = root_path
    @file_path = file_path
    @service_name = service_name
    @package_simple_name = File.basename(@file_path, ".rb")
    @package_name = get_package_name
    @initialized = true
  end

  def get_package_name
    root = Pathname.new(@root_path)
    file = Pathname.new(File.dirname(@file_path))
    relative_path = file.relative_path_from(root)
    relative_path = relative_path.to_s.gsub(File::SEPARATOR, PACKAGE_SEPARATOR)
    if relative_path == '.'
      @package_simple_name
    else
      [relative_path, @package_simple_name].join(PACKAGE_SEPARATOR)
    end
  end

  def get_full_name(simple_name)
    if @current_parent
      full_name = [@current_parent, simple_name].compact.join(PARENT_SEPARATOR)
    else
      full_name = [@package_name, simple_name].compact.join(PACKAGE_NAME_SEPARATOR)
    end
    @current_parent = full_name
    full_name
  end

  def on_module(node)
    previous_parent = @current_parent
    full_name = get_full_name(node.identifier.const_name)
    @logger.debug("Module: #{full_name} (line: #{node.loc.line})")
    super
    @current_parent = previous_parent
  end

  def on_class(node)
    if @depth > MAX_CLASS_DEPTH
      @logger.debug("Skipping class #{node.identifier.const_name} (package: #{@package_name})")
      return
    end
    previous_parent = @current_parent
    previous_class = @last_class
    simple_name = node.identifier.const_name
    full_name = get_full_name(simple_name)
    @last_class = full_name
    content = node.source
    object_ = Object_.new(simple_name, full_name, @file_path, @service_name, content)
    @objects_ << object_
    begin
      if node.parent_class
        @logger.debug("Class: #{full_name} from parent #{node.parent_class.const_name} (line: #{node.loc.line}:#{node.loc.column})")
      else
        @logger.debug("Class: #{full_name} (line: #{node.loc.line}:#{node.loc.column})")
      end
    rescue NoMethodError => e
      @logger.error("Error in node : #{node}")
      @logger.error("Error: #{e.message}")
    end
    @depth += 1
    super
    @depth -= 1
    @current_parent = previous_parent
    @last_class = previous_class
  end

  def on_def(node)
    if (@depth + @method_depth) > MAX_METHOD_DEPTH
      @logger.debug("Skipping method #{node.method_name} (package: #{@package_name})")
      return
    end
    previous_parent = @current_parent
    simple_name = "#{node.method_name}()"
    full_name = get_full_name(simple_name)
    content = node.source
    parent_name = ""
    if @last_class
      parent_name = @last_class
    end
    executable_ = Executable_.new(simple_name, full_name, parent_name, @service_name, content)
    @executables_ << executable_
    if node.class_constructor?
      @logger.debug("Constructor: #{full_name} (line: #{node.loc.line})")
    else
      @logger.debug("Method: #{full_name} (line: #{node.loc.line})")
    end
    @depth += 1
    @method_depth += 1
    super
    @method_depth -= 1
    @depth -= 1
    @current_parent = previous_parent
  end

  def process(node)
    if @initialized
      @logger.debug("Traversing nodes of #{@package_name}")
      super
    else
      @logger.error("Processor not initialized. Call init_vars before processing.")
      throw "Processor not initialized"
    end
  end
end