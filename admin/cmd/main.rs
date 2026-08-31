package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}package main

import (
	"flag"
	"fmt"
	"os"

	"admin/internal/services/auth"
)

func create_admin() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=...\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")

		err := adminCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}
		auth.CreateAdminUser(*adminUsername, *adminPassword)
		
		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)

	default:
		flag.Usage()
		os.Exit(1)
	}
}