package workspace

type Artifact struct {
	Name string
	Size int64
}

func CaptureArtifact(name string, size int64) Artifact {
	return Artifact{Name: name, Size: size}
}
