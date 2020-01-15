// +build aix js nacl plan9 windows android solaris

package main

func getwinsize() winsize {
	return winsize{80, 40, 100, 100} // complete fudge
}
