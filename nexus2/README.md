# Testing Requirements

This is a WIP to be coded in go. For the meantime, do the following to create a test data.

```bash
# Create a file
touch resources/artifactA.jar

# Upload file to nexus. This requires maven installed!
mvn deploy:deploy-file -DgroupId=com.example -DartifactId=artifactA -Dversion=1.0.0 -DgeneratePom=true -Dpackage=jar -DrepositoryId=releases -Durl=http://admin:admin123@localhost:8081/nexus/content/repositories/releases/ -Dfile=resources/artifactA.jar
```