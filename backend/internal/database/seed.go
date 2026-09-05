package database

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jaberpatwary/startech/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Database already seeded, skipping.")
		return nil
	}

	log.Println("Seeding database...")

	// --- Admin User ---
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		ID:           uuid.NewString(),
		Name:         "MI-Tech Admin",
		Email:        "admin@mitech.com",
		Phone:        "01700000001",
		PasswordHash: string(hash),
		Role:         models.RoleAdmin,
	}
	db.Create(&admin)

	// --- Categories ---
	cats := []struct{ name, icon string }{
		{"Desktop", "monitor"},
		{"Laptop", "laptop"},
		{"Component", "cpu"},
		{"Monitor", "display"},
	}
	catMap := map[string]string{}
	for i, c := range cats {
		cat := models.Category{
			ID:        uuid.NewString(),
			Name:      c.name,
			Slug:      slug.Make(c.name),
			Icon:      c.icon,
			SortOrder: i,
		}
		db.Create(&cat)
		catMap[c.name] = cat.ID
	}

	// --- Brands ---
	brands := []string{"ASUS", "MSI", "Gigabyte", "Intel", "AMD", "Corsair", "Samsung", "Dell"}
	brandMap := map[string]string{}
	for _, b := range brands {
		brand := models.Brand{
			ID:   uuid.NewString(),
			Name: b,
			Slug: slug.Make(b),
		}
		db.Create(&brand)
		brandMap[b] = brand.ID
	}

	// --- Products ---
	type productSeed struct {
		name, sku, category, brand, short, desc string
		price, discount, stock                  int
		featured                                bool
		specs                                   map[string]string
		images                                  []string
	}

	demoImg := func(keyword string) []string {
		return []string{
			"https://placehold.co/800x600/1a1a2e/ef4a23?text=" + keyword,
			"https://placehold.co/800x600/16213e/ef4a23?text=" + keyword + "+2",
		}
	}

	products := []productSeed{
		// Laptops
		{name: "ASUS ROG Strix G16 Gaming Laptop", sku: "ASUS-ROG-G16-2024", category: "Laptop", brand: "ASUS",
			short: "Intel Core i7-14650HX, RTX 4060 8GB, 16GB DDR5, 1TB NVMe SSD, 165Hz Display",
			desc:  "Power up your game with the ROG Strix G16 featuring 14th Gen Intel Core and NVIDIA GeForce RTX 40 Series GPU. Extreme gaming performance meets stunning visuals.",
			price: 185000, discount: 175000, stock: 12, featured: true,
			specs:  map[string]string{"Processor": "Intel Core i7-14650HX", "RAM": "16GB DDR5", "Storage": "1TB NVMe SSD", "GPU": "NVIDIA RTX 4060 8GB", "Display": "16\" 165Hz FHD", "OS": "Windows 11 Home"},
			images: demoImg("ASUS+ROG")},
		{name: "Dell XPS 15 Ultrabook", sku: "DELL-XPS15-2024", category: "Laptop", brand: "Dell",
			short: "Intel Core i9-14900H, NVIDIA RTX 4070 8GB, 32GB DDR5, 1TB SSD, OLED 4K",
			desc:  "Dell XPS 15 delivers extraordinary performance in a sleek, ultra-thin design. Premium OLED display with vivid colors for professionals.",
			price: 245000, stock: 8, featured: true,
			specs:  map[string]string{"Processor": "Intel Core i9-14900H", "RAM": "32GB DDR5", "Storage": "1TB SSD", "GPU": "RTX 4070 8GB", "Display": "15.6\" 4K OLED", "Battery": "86Wh"},
			images: demoImg("Dell+XPS")},
		{name: "MSI Stealth 16 AI Studio", sku: "MSI-STEALTH16-2024", category: "Laptop", brand: "MSI",
			short: "Intel Core Ultra 9 185H, RTX 4080 12GB, 64GB DDR5, 2TB SSD, 4K OLED",
			desc:  "The ultimate creator and gaming laptop with top-tier specs. RTX 4080 Laptop GPU brings desktop-class performance wherever you go.",
			price: 320000, discount: 310000, stock: 5, featured: true,
			specs:  map[string]string{"Processor": "Intel Core Ultra 9 185H", "RAM": "64GB DDR5", "Storage": "2TB SSD", "GPU": "RTX 4080 12GB", "Display": "16\" 4K OLED 120Hz"},
			images: demoImg("MSI+Stealth")},

		// Desktops
		{name: "ASUS ROG Strix GA35 Gaming Desktop", sku: "ASUS-ROG-GA35-PC", category: "Desktop", brand: "ASUS",
			short: "AMD Ryzen 9 7950X, RTX 4090 24GB, 64GB DDR5, 2TB NVMe SSD",
			desc:  "The ROG Strix GA35 is a beast of a desktop PC, built to crush any game at any resolution. Insane performance for professional creators too.",
			price: 480000, stock: 6, featured: true,
			specs:  map[string]string{"Processor": "AMD Ryzen 9 7950X", "RAM": "64GB DDR5", "Storage": "2TB NVMe", "GPU": "NVIDIA RTX 4090 24GB", "OS": "Windows 11 Pro"},
			images: demoImg("ASUS+Desktop")},
		{name: "MSI Infinite RS Gaming PC", sku: "MSI-INFINITE-RS", category: "Desktop", brand: "MSI",
			short: "Intel Core i9-14900K, RTX 4080 Super 16GB, 32GB DDR5, 2TB SSD",
			desc:  "MSI Infinite RS delivers raw gaming power with an Intel Core i9 and an RTX 4080 Super, ready for 4K gaming and content creation.",
			price: 380000, discount: 365000, stock: 4, featured: false,
			specs:  map[string]string{"Processor": "Intel Core i9-14900K", "RAM": "32GB DDR5", "GPU": "RTX 4080 Super 16GB", "Storage": "2TB NVMe SSD"},
			images: demoImg("MSI+Desktop")},

		// Components
		{name: "Intel Core i9-14900K Processor", sku: "INTEL-I9-14900K", category: "Component", brand: "Intel",
			short: "24-Core (8P+16E), Max 6.0GHz, 36MB Cache, LGA1700, No Cooling",
			desc:  "The Intel Core i9-14900K is Intel's flagship desktop processor with 24 cores. Phenomenal performance for gaming, rendering and multitasking.",
			price: 72000, discount: 68000, stock: 25, featured: true,
			specs:  map[string]string{"Cores": "24 (8P+16E)", "Max Turbo": "6.0 GHz", "Cache": "36MB L3", "Socket": "LGA1700", "TDP": "125W Base / 253W MTP"},
			images: demoImg("Intel+i9")},
		{name: "AMD Ryzen 9 7950X3D Processor", sku: "AMD-R9-7950X3D", category: "Component", brand: "AMD",
			short: "16-Core 32-Thread, 5.7GHz Max Boost, 144MB Cache, AM5 Socket",
			desc:  "AMD's Ryzen 9 7950X3D combines high-core-count Zen 4 architecture with AMD 3D V-Cache for unmatched gaming and productivity performance.",
			price: 85000, stock: 15, featured: true,
			specs:  map[string]string{"Cores": "16 (32 Threads)", "Max Boost": "5.7 GHz", "Cache": "144MB Total", "Socket": "AM5", "TDP": "120W"},
			images: demoImg("AMD+Ryzen")},
		{name: "Corsair Vengeance 32GB DDR5-6000 RGB", sku: "CORSAIR-VEN-32GB-DDR5", category: "Component", brand: "Corsair",
			short: "2x 16GB DDR5-6000MHz CL36, RGB, Intel XMP 3.0 / AMD EXPO",
			desc:  "Corsair Vengeance DDR5 is fast, reliable memory designed for Intel 12th/13th/14th Gen and AMD Ryzen 7000 series platforms.",
			price: 18500, discount: 17000, stock: 40, featured: false,
			specs:  map[string]string{"Capacity": "32GB (2x16GB)", "Speed": "6000MHz DDR5", "Latency": "CL36", "Voltage": "1.35V", "RGB": "Yes"},
			images: demoImg("Corsair+RAM")},
		{name: "Gigabyte RTX 4080 Super Gaming OC 16GB", sku: "GIGABYTE-RTX4080S-OC", category: "Component", brand: "Gigabyte",
			short: "16GB GDDR6X, Boost Clock 2595MHz, WINDFORCE 3X Cooling",
			desc:  "The Gigabyte RTX 4080 Super Gaming OC brings extreme 4K gaming performance with WINDFORCE 3X cooling technology.",
			price: 150000, discount: 142000, stock: 10, featured: true,
			specs:  map[string]string{"VRAM": "16GB GDDR6X", "Boost Clock": "2595 MHz", "Cooling": "WINDFORCE 3X", "Power": "320W TDP", "Outputs": "3x DP 1.4 + HDMI 2.1"},
			images: demoImg("Gigabyte+RTX")},
		{name: "Gigabyte Z790 Aorus Elite AX Motherboard", sku: "GIGABYTE-Z790-ELITE", category: "Component", brand: "Gigabyte",
			short: "LGA1700, DDR5, PCIe 5.0, WiFi 6E, Thunderbolt 4, ATX Form Factor",
			desc:  "The Gigabyte Z790 Aorus Elite AX supports Intel 12th/13th/14th Gen processors with DDR5 memory and PCIe 5.0 for future-proof performance.",
			price: 32000, stock: 18, featured: false,
			specs:  map[string]string{"Socket": "LGA1700", "Chipset": "Intel Z790", "Memory": "DDR5 up to 7600MHz", "Form Factor": "ATX", "Network": "2.5G LAN + WiFi 6E"},
			images: demoImg("Gigabyte+MB")},
		{name: "Corsair MP700 Pro 2TB NVMe SSD", sku: "CORSAIR-MP700-2TB", category: "Component", brand: "Corsair",
			short: "PCIe Gen5, 12,400 MB/s Read, 11,800 MB/s Write, M.2 2280",
			desc:  "The fastest SSD for enthusiasts. Corsair MP700 Pro Gen5 delivers mind-blowing transfer speeds for demanding workloads.",
			price: 28000, discount: 25000, stock: 20, featured: false,
			specs:  map[string]string{"Capacity": "2TB", "Interface": "PCIe Gen5 NVMe", "Read Speed": "12,400 MB/s", "Write Speed": "11,800 MB/s", "Form Factor": "M.2 2280"},
			images: demoImg("Corsair+SSD")},

		// Monitors
		{name: "Samsung Odyssey G9 49\" Super Ultrawide", sku: "SAMSUNG-OD-G9-49", category: "Monitor", brand: "Samsung",
			short: "49\" Dual QHD 5120x1440, 240Hz, 1ms, VA QLED, Curved 1800R",
			desc:  "The ultimate gaming monitor. Samsung Odyssey G9 provides the full dual monitor experience in one epic 49-inch QLED panel.",
			price: 185000, discount: 175000, stock: 8, featured: true,
			specs:  map[string]string{"Size": "49 inches", "Resolution": "5120x1440 (DQHD)", "Refresh Rate": "240Hz", "Response Time": "1ms GTG", "Panel": "VA QLED", "Curvature": "1800R"},
			images: demoImg("Samsung+Odyssey")},
		{name: "ASUS ROG Swift PG27AQN 27\" Gaming Monitor", sku: "ASUS-ROG-PG27AQN", category: "Monitor", brand: "ASUS",
			short: "27\" 2560x1440 QHD, 360Hz, 1ms, IPS, G-Sync Ultimate, HDR600",
			desc:  "ASUS ROG Swift PG27AQN is designed for elite esports players. 360Hz refresh rate with QHD resolution for crisp, fluid gameplay.",
			price: 95000, stock: 15, featured: true,
			specs:  map[string]string{"Size": "27 inches", "Resolution": "2560x1440 QHD", "Refresh Rate": "360Hz", "Response": "1ms GTG", "HDR": "HDR600", "Sync": "G-Sync Ultimate"},
			images: demoImg("ASUS+Monitor")},
		{name: "Dell UltraSharp U3223QE 32\" 4K Monitor", sku: "DELL-U3223QE-32", category: "Monitor", brand: "Dell",
			short: "32\" 4K UHD 3840x2160, IPS Black, USB-C 96W, 60Hz, for Professionals",
			desc:  "Dell UltraSharp 32 4K USB-C Hub Monitor with IPS Black technology for exceptional contrast and color accuracy.",
			price: 68000, discount: 64000, stock: 12, featured: false,
			specs:  map[string]string{"Size": "32 inches", "Resolution": "3840x2160 (4K UHD)", "Panel": "IPS Black", "Color": "99% sRGB, 98% DCI-P3", "USB-C": "96W Power Delivery"},
			images: demoImg("Dell+Monitor")},
		{name: "MSI MAG 341CQP QD-OLED 34\" Monitor", sku: "MSI-MAG-341CQP", category: "Monitor", brand: "MSI",
			short: "34\" 3440x1440 QD-OLED, 175Hz, 0.1ms, Curved 1800R, HDR True Black 400",
			desc:  "Immerse yourself in vivid colors with the MSI MAG 341CQP QD-OLED ultrawide monitor. Perfect balance of gaming and creative work.",
			price: 78000, stock: 9, featured: false,
			specs:  map[string]string{"Size": "34 inches", "Resolution": "3440x1440 UWQHD", "Panel": "QD-OLED", "Refresh": "175Hz", "Response": "0.1ms", "HDR": "True Black 400"},
			images: demoImg("MSI+Monitor")},
		// Extra products to reach 20+
		{name: "AMD Radeon RX 7900 XTX 24GB", sku: "AMD-RX7900XTX-24GB", category: "Component", brand: "AMD",
			short: "24GB GDDR6, 2500MHz Boost, 384-bit Memory, 355W TDP",
			desc:  "AMD's fastest Radeon GPU ever. The RX 7900 XTX delivers exceptional 4K gaming performance and ray tracing capabilities.",
			price: 125000, stock: 7, featured: false,
			specs:  map[string]string{"VRAM": "24GB GDDR6", "Boost Clock": "2500 MHz", "Memory Bus": "384-bit", "TDP": "355W", "Outputs": "DP 2.1 + HDMI 2.1"},
			images: demoImg("AMD+RX7900")},
		{name: "Corsair RM1000x 1000W 80+ Gold PSU", sku: "CORSAIR-RM1000X-1000W", category: "Component", brand: "Corsair",
			short: "1000W, 80 Plus Gold, Fully Modular, Zero RPM Fan Mode, 10 Year Warranty",
			desc:  "The Corsair RM1000x provides reliable, efficient power with fully modular cables and an ultra-quiet zero RPM fan mode.",
			price: 22000, discount: 20000, stock: 30, featured: false,
			specs:  map[string]string{"Wattage": "1000W", "Efficiency": "80 Plus Gold", "Modular": "Fully Modular", "Fan": "Zero RPM Mode", "Warranty": "10 Years"},
			images: demoImg("Corsair+PSU")},
		{name: "Samsung 990 Pro 2TB NVMe SSD", sku: "SAMSUNG-990PRO-2TB", category: "Component", brand: "Samsung",
			short: "PCIe Gen4, 7,450 MB/s Read, M.2 2280, MLC NAND, RGB",
			desc:  "Samsung 990 Pro offers blazing fast PCIe 4.0 speeds with improved endurance and efficiency for gaming and professional use.",
			price: 19000, stock: 35, featured: false,
			specs:  map[string]string{"Capacity": "2TB", "Interface": "PCIe Gen4 x4 NVMe", "Read Speed": "7,450 MB/s", "Write Speed": "6,900 MB/s"},
			images: demoImg("Samsung+SSD")},
		{name: "ASUS TUF Gaming F15 Laptop", sku: "ASUS-TUF-F15-2024", category: "Laptop", brand: "ASUS",
			short: "Intel Core i7-13700H, RTX 4070 8GB, 16GB DDR5, 512GB SSD, 144Hz",
			desc:  "ASUS TUF Gaming F15 is built to last with MIL-STD-810H military-grade durability, offering great gaming value.",
			price: 135000, discount: 128000, stock: 18, featured: false,
			specs:  map[string]string{"Processor": "Intel Core i7-13700H", "GPU": "RTX 4070 8GB", "RAM": "16GB DDR5", "Storage": "512GB SSD", "Display": "15.6\" 144Hz FHD"},
			images: demoImg("ASUS+TUF")},
		{name: "MSI PRO DP180 Mini PC", sku: "MSI-PRO-DP180", category: "Desktop", brand: "MSI",
			short: "Intel Core i5-14400F, 16GB DDR5, 512GB SSD, Compact Form Factor",
			desc:  "Powerful mini PC for home office and light gaming. MSI PRO DP180 packs desktop performance into a tiny, elegant package.",
			price: 65000, stock: 10, featured: false,
			specs:  map[string]string{"Processor": "Intel Core i5-14400F", "RAM": "16GB DDR5", "Storage": "512GB NVMe", "Form Factor": "Mini PC"},
			images: demoImg("MSI+MiniPC")},
		{name: "Samsung 32\" ViewFinity S8 4K Monitor", sku: "SAMSUNG-S8-32-4K", category: "Monitor", brand: "Samsung",
			short: "32\" 4K UHD 3840x2160, 60Hz, IPS, 99% sRGB, USB-C 65W, HDR400",
			desc:  "Samsung ViewFinity S8 provides a spacious 4K workspace for professionals who demand color accuracy and productivity.",
			price: 45000, discount: 42000, stock: 20, featured: false,
			specs:  map[string]string{"Size": "32 inches", "Resolution": "3840x2160 4K", "Color": "99% sRGB", "USB-C": "65W PD", "HDR": "HDR400"},
			images: demoImg("Samsung+4K")},
	}

	for _, p := range products {
		discountPrice := (*int)(nil)
		if p.discount > 0 {
			discountPrice = &p.discount
		}

		specsJSON, _ := json.Marshal(p.specs)

		product := models.Product{
			ID:               uuid.NewString(),
			Name:             p.name,
			Slug:             slug.Make(p.name),
			SKU:              p.sku,
			CategoryID:       catMap[p.category],
			BrandID:          brandMap[p.brand],
			Price:            p.price,
			DiscountPrice:    discountPrice,
			Stock:            p.stock,
			ShortDescription: p.short,
			Description:      p.desc,
			Specs:            datatypes.JSON(specsJSON),
			Status:           models.ProductActive,
			Featured:         p.featured,
		}
		db.Create(&product)

		for i, imgURL := range p.images {
			db.Create(&models.ProductImage{
				ID:        uuid.NewString(),
				ProductID: product.ID,
				URL:       imgURL,
				SortOrder: i,
			})
		}
	}

	// --- Banners ---
	banners := []models.Banner{
		{ID: uuid.NewString(), Title: "Next-Gen Gaming Laptops", Subtitle: "Up to ৳15,000 Off on ROG & MSI Gaming Laptops", Image: "https://placehold.co/1400x500/0f172a/ef4a23?text=Gaming+Laptops+Sale", Link: "/category/laptop", SortOrder: 0, Active: true},
		{ID: uuid.NewString(), Title: "Build Your Dream PC", Subtitle: "Latest Intel & AMD Processors — In Stock Now", Image: "https://placehold.co/1400x500/1e1b4b/ef4a23?text=PC+Components", Link: "/category/component", SortOrder: 1, Active: true},
		{ID: uuid.NewString(), Title: "Immersive Monitors", Subtitle: "4K & QD-OLED Gaming Monitors from MI-Tech", Image: "https://placehold.co/1400x500/052e16/ef4a23?text=Gaming+Monitors", Link: "/category/monitor", SortOrder: 2, Active: true},
	}
	for _, b := range banners {
		db.Create(&b)
	}

	// --- Coupons ---
	coupons := []models.Coupon{
		{ID: uuid.NewString(), Code: "MITECH10", Type: models.CouponPercent, Value: 10, MinOrder: 10000, MaxUses: 100, Active: true},
		{ID: uuid.NewString(), Code: "WELCOME500", Type: models.CouponFixed, Value: 500, MinOrder: 5000, MaxUses: 200, Active: true},
	}
	for _, c := range coupons {
		db.Create(&c)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}
