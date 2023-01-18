package main

import (
	"errors"
	"fmt"
	"github.com/bittorrent/btfs-encryption/btfs"
	"github.com/bittorrent/btfs-encryption/enc"
	"github.com/urfave/cli/v2"
	"os"
	"path"
)

var cmds = &cli.App{
	Name:    "btfs-encryption",
	Version: "v0.1.0",
	Usage:   "btfs-encryption is a demo project for btfs encryption protocol",
	Commands: []*cli.Command{
		{
			Name:  "encrypt",
			Usage: "Encrypt local file or folder and add it to BTFS",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "pub", Value: path.Join(os.Getenv("HOME"), "/.ssh/id_rsa.pub"), Usage: "encrypt public key file path"},
			},
			ArgsUsage: "source_file_or_folder_path",
			Before: func(ctx *cli.Context) error {
				return btfs.Init()
			},
			Action: func(ctx *cli.Context) (err error) {
				if ctx.NArg() < 1 {
					err = errors.New("arguments not enough")
					return
				}
				srcPath := path.Clean(ctx.Args().Get(0))
				pubPath := path.Clean(ctx.String("pub"))
				rst, err := enc.EncryptToBTFS(srcPath, pubPath)
				if err != nil {
					return
				}
				fmt.Printf("Encrypted File: \nCID  - %s\nName - %s\nSize - %s\n", rst.Hash, rst.Name, rst.Size)
				return
			},
		},
		{
			Name:  "decrypt",
			Usage: "Get encrypted file from BTFS and decrypt it to local",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dst", Value: "./", Usage: "destination directory path"},
				&cli.StringFlag{Name: "prv", Value: path.Join(os.Getenv("HOME"), "/.ssh/id_rsa.pub"), Usage: "private key path"},
			},
			ArgsUsage: "cid",
			Before: func(ctx *cli.Context) error {
				return btfs.Init()
			},
			Action: func(ctx *cli.Context) (err error) {
				if ctx.NArg() < 1 {
					err = errors.New("arguments not enough")
					return
				}
				cid := ctx.Args().Get(0)
				dstPath := path.Clean(ctx.String("dst"))
				prvPath := path.Clean(ctx.String("pub"))
				err = enc.DecryptFromBTFS(cid, dstPath, prvPath)
				if err != nil {
					return
				}
				fmt.Println("completed!")
				return
			},
		},
		{
			Name:  "encrypt-local",
			Usage: "Encrypt local file or folder",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dst", Value: "./", Usage: "destination directory path"},
				&cli.StringFlag{Name: "pub", Value: path.Join(os.Getenv("HOME"), "/.ssh/id_rsa.pub"), Usage: "encrypt public key file path"},
			},
			ArgsUsage: "source_file_or_folder_path",
			Action: func(ctx *cli.Context) (err error) {
				if ctx.NArg() < 1 {
					err = errors.New("arguments not enough")
					return
				}
				srcPath := path.Clean(ctx.Args().Get(0))
				dstPath := path.Clean(ctx.String("dst"))
				pubPath := path.Clean(ctx.String("pub"))
				err = enc.EncryptToLocal(srcPath, dstPath, pubPath)
				if err != nil {
					return
				}
				fmt.Println("completed!")
				return
			},
		},
		{
			Name:  "decrypt-local",
			Usage: "Decrypt local encrypted file",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dst", Value: "./", Usage: "destination directory path"},
				&cli.StringFlag{Name: "prv", Value: path.Join(os.Getenv("HOME"), "/.ssh/id_rsa.pub"), Usage: "private key path"},
			},
			ArgsUsage: "source_file_path",
			Before: func(ctx *cli.Context) error {
				return btfs.Init()
			},
			Action: func(ctx *cli.Context) (err error) {
				if ctx.NArg() < 1 {
					err = errors.New("arguments not enough")
					return
				}
				srcPath := path.Clean(ctx.Args().Get(0))
				dstPath := path.Clean(ctx.String("dst"))
				prvPath := path.Clean(ctx.String("pub"))
				err = enc.DecryptFromLocal(srcPath, dstPath, prvPath)
				if err != nil {
					return
				}
				fmt.Println("completed!")
				return
			},
		},
	},
}
