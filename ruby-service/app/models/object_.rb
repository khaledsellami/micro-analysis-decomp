# frozen_string_literal: true

class Object_
  attr_accessor :simpleName, :fullName, :filePath, :serviceName, :content, :isInterface, :isAnnotation
  def initialize(simpleName, fullName, filePath, serviceName, content)
    @simpleName = simpleName
    @fullName = fullName
    @filePath = filePath
    @serviceName = serviceName
    @content = content
    @isInterface = false
    @isAnnotation = false
  end

  def to_hash
    {
      'simpleName' => @simpleName,
      'fullName' => @fullName,
      'filePath' => @filePath,
      'serviceName' => @serviceName,
      'content' => @content,
      'isInterface' => @isInterface,
      'isAnnotation' => @isAnnotation
    }
  end
end
