# frozen_string_literal: true

class Executable_
  attr_accessor :simpleName, :fullName, :parentName, :serviceName, :content
  def initialize(simpleName, fullName, parentName, serviceName, content)
    @simpleName = simpleName
    @fullName = fullName
    @parentName = parentName
    @serviceName = serviceName
    @content = content
  end

  def to_hash
    {
      'simpleName' => @simpleName,
      'fullName' => @fullName,
      'parentName' => @parentName,
      'serviceName' => @serviceName,
      'content' => @content
    }
  end
end
