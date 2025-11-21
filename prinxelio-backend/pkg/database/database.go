package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
)

// InitDB initializes the database connection and creates the tables.
func InitDB() (*sql.DB, error) {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the database!")

	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		phone_number VARCHAR(20) NOT NULL UNIQUE,
		otp_code VARCHAR(6) NULL,
		otp_created_at TIMESTAMP NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP NULL,
		status TINYINT(1) DEFAULT 1
	);`

	createProductTable := `
	CREATE TABLE IF NOT EXISTS product (
		id INT AUTO_INCREMENT PRIMARY KEY,
		product_name VARCHAR(255) NOT NULL,
		product_image TEXT NULL,
		product_price DECIMAL(15,2) NOT NULL,
		product_discount DECIMAL(15,2) DEFAULT 0,
		product_discount_amount INT DEFAULT 0,
		product_desc TEXT,
		product_viewed INT DEFAULT 0,
		product_downloaded INT DEFAULT 0,
		product_create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

    createTransactionsTable := `
    CREATE TABLE IF NOT EXISTS transactions (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT NOT NULL,
        product_id INT NOT NULL,
        base_price DECIMAL(15,2) NOT NULL,
        admin_fee DECIMAL(15,2) DEFAULT 0,
        total_amount DECIMAL(15,2) NOT NULL,
        status ENUM('UNPAID', 'PAID', 'FAILED', 'EXPIRED') DEFAULT 'UNPAID',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        merchant_ref VARCHAR(100) NULL,
        expired_time INT NULL,
        reference VARCHAR(100) UNIQUE NOT NULL,
        link_qr TEXT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (product_id) REFERENCES product(id)
    );`

    createCategoryTable := `
    CREATE TABLE IF NOT EXISTS category (
        id INT AUTO_INCREMENT PRIMARY KEY,
        category_name VARCHAR(255) NOT NULL,
        category_images TEXT NULL,
        category_color VARCHAR(7) NULL,
        category_create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

    _, err := db.Exec(createUserTable)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

    _, err = db.Exec(createProductTable)
    if err != nil {
        return fmt.Errorf("error creating product table: %v", err)
    }

    _, err = db.Exec(createTransactionsTable)
	if err != nil {
		return fmt.Errorf("error creating transactions table: %v", err)
	}

    _, err = db.Exec(createCategoryTable)
    if err != nil {
        return fmt.Errorf("error creating category table: %v", err)
    }

    var catImgCol int
    _ = db.QueryRow("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = ? AND table_name = 'category' AND column_name = 'category_images'", os.Getenv("DB_NAME")).Scan(&catImgCol)
    if catImgCol == 0 {
        _, _ = db.Exec("ALTER TABLE category ADD COLUMN category_images TEXT NULL")
    }

    var catColorCol int
    _ = db.QueryRow("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = ? AND table_name = 'category' AND column_name = 'category_color'", os.Getenv("DB_NAME")).Scan(&catColorCol)
    if catColorCol == 0 {
        _, _ = db.Exec("ALTER TABLE category ADD COLUMN category_color VARCHAR(7) NULL")
    }

    var colExists int
    _ = db.QueryRow("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = ? AND table_name = 'product' AND column_name = 'product_category'", os.Getenv("DB_NAME")).Scan(&colExists)
    if colExists == 0 {
        _, _ = db.Exec("ALTER TABLE product ADD COLUMN product_category INT NULL")
    }

    var fkExists int
    _ = db.QueryRow("SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE WHERE table_schema = ? AND table_name = 'product' AND column_name = 'product_category' AND referenced_table_name = 'category'", os.Getenv("DB_NAME")).Scan(&fkExists)
    if fkExists == 0 {
        _, _ = db.Exec("ALTER TABLE product ADD CONSTRAINT fk_product_category FOREIGN KEY (product_category) REFERENCES category(id)")
    }

    _, _ = db.Exec("ALTER TABLE transactions MODIFY COLUMN status ENUM('UNPAID','PAID','FAILED','EXPIRED') DEFAULT 'UNPAID'")
    fmt.Println("Tables created successfully!")
    return nil
}

// SeedProducts adds some sample products to the database.
func SeedProducts(db *sql.DB) error {
	// Check if products already exist to avoid duplication
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM product").Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking product count: %v", err)
	}

	if count > 0 {
		fmt.Println("Products table is not empty, skipping seeding.")
		return nil
	}

    products := []struct {
        Name        string
        Image       string
        Price       float64
        Discount    float64
        DiscountAmt int
        Desc        string
        CatName     string
    }{
        {
            Name:        "Ebook Panduan Go Lengkap",
            Image:       "https://via.placeholder.com/400x300.png/000000/FFFFFF?text=Ebook+Go",
            Price:       150000,
            Discount:    120000,
            DiscountAmt: 20,
            Desc:        "Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.",
            CatName:     "Ebook",
        },
        {
            Name:        "Template Website Portofolio",
            Image:       "https://via.placeholder.com/400x300.png/000000/FFFFFF?text=Template+Web",
            Price:       75000,
            Discount:    0,
            DiscountAmt: 0,
            Desc:        "Template website portofolio modern dan responsif menggunakan HTML, CSS, dan JavaScript.",
            CatName:     "Website",
        },
        {
            Name:        "Koleksi Ikon Vektor Premium",
            Image:       "https://via.placeholder.com/400x300.png/000000/FFFFFF?text=Ikon+Vektor",
            Price:       200000,
            Discount:    150000,
            DiscountAmt: 25,
            Desc:        "Lebih dari 1000 ikon vektor premium untuk kebutuhan desain Anda.",
            CatName:     "Desain",
        },
    }

    query := `INSERT INTO product (product_name, product_image, product_price, product_discount, product_discount_amount, product_desc, product_category) VALUES (?, ?, ?, ?, ?, ?, ?)`

    var websiteID sql.NullInt64
    _ = db.QueryRow("SELECT id FROM category WHERE category_name = ?", "Website").Scan(&websiteID)

    for _, p := range products {
        var catID sql.NullInt64
        _ = db.QueryRow("SELECT id FROM category WHERE category_name = ?", p.CatName).Scan(&catID)
        var cid interface{}
        if catID.Valid {
            cid = catID.Int64
        } else if strings.EqualFold(p.CatName, "Website") && websiteID.Valid {
            cid = websiteID.Int64
        } else {
            cid = nil
        }
        _, err := db.Exec(query, p.Name, p.Image, p.Price, p.Discount, p.DiscountAmt, p.Desc, cid)
        if err != nil {
            if !isDuplicateEntryError(err) {
                return fmt.Errorf("error seeding product %s: %v", p.Name, err)
            }
        }
    }

	fmt.Println("Products seeded successfully!")
    return nil
}

// SeedCategories adds sample categories to the database.
func SeedCategories(db *sql.DB) error {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM category").Scan(&count)
    if err != nil {
        return fmt.Errorf("error checking category count: %v", err)
    }
    if count > 0 {
        fmt.Println("Category table is not empty, skipping seeding.")
        return nil
    }
    cats := []struct{
        Name string
        Image string
        Color string
        Created string
    }{
        {Name:"Website", Image:"https://example.com/cat1.jpg", Color:"#FF5733", Created:"2023-01-01"},
        {Name:"Ebook", Image:"https://example.com/cat2.jpg", Color:"#27AE60", Created:"2023-01-02"},
        {Name:"Desain", Image:"https://example.com/cat3.jpg", Color:"#2980B9", Created:"2023-01-03"},
        {Name:"AI Tools", Image:"https://example.com/cat4.jpg", Color:"#8E44AD", Created:"2023-01-04"},
        {Name:"Audio", Image:"https://example.com/cat5.jpg", Color:"#F39C12", Created:"2023-01-05"},
    }
    for _, c := range cats {
        _, err := db.Exec("INSERT INTO category (category_name, category_images, category_color, category_create_at) VALUES (?, ?, ?, ?)", c.Name, c.Image, c.Color, c.Created)
        if err != nil {
            if !isDuplicateEntryError(err) {
                return fmt.Errorf("error seeding category %s: %v", c.Name, err)
            }
        }
    }
    fmt.Println("Categories seeded successfully!")
    return nil
}

// isDuplicateEntryError checks if an error is a MySQL duplicate entry error.
func isDuplicateEntryError(err error) bool {
    mysqlErr, ok := err.(*mysql.MySQLError)
    return ok && mysqlErr.Number == 1062
}
