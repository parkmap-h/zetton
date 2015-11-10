lazy val deployStaging = taskKey[Unit]("deploy staging")

deployStaging := {
  "goapp deploy -application parkmap-h-staging src" !
}

lazy val deploy = taskKey[Unit]("deploy")

deploy := {
  "goapp deploy -application parkmap-h src" !
}
