//go:build !windows

package core

func launchUpdateWorker(target,source,digest,tag,old string)error{return errUpdateWorkerUnsupported}
func updateWorkerCLI(args []string){die("%v",errUpdateWorkerUnsupported)}
