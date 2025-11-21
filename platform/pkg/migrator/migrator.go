package migrator

type IMigrator interface {
	Up() error
}
