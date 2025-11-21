-- phpMyAdmin SQL Dump
-- version 5.2.2
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Waktu pembuatan: 21 Nov 2025 pada 08.27
-- Versi server: 11.4.8-MariaDB-ubu2404
-- Versi PHP: 8.3.25

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `user_prinxelio_digital`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `category`
--

CREATE TABLE `category` (
  `id` int(11) NOT NULL,
  `category_name` varchar(255) NOT NULL,
  `category_create_at` timestamp NULL DEFAULT current_timestamp(),
  `category_images` text DEFAULT NULL,
  `category_color` varchar(7) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;

--
-- Dumping data untuk tabel `category`
--

INSERT INTO `category` (`id`, `category_name`, `category_create_at`, `category_images`, `category_color`) VALUES
(1, 'WEBSITE', '2022-12-31 17:00:00', '/images/category/website.png', '#000000'),
(2, 'N8N', '2023-01-01 17:00:00', '/images/category/n8n.png', '#5c002e'),
(3, 'DOCS', '2023-01-02 17:00:00', '/images/category/docs.png', '#2c005c'),
(4, 'AI TOOLS', '2023-01-03 17:00:00', '/images/category/ai.png', '#00bfc9'),
(5, 'IMAGES', '2023-01-04 17:00:00', '/images/category/images.png', '#3b1f00');

-- --------------------------------------------------------

--
-- Struktur dari tabel `logs_file_transfer`
--

CREATE TABLE `logs_file_transfer` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `transaction_id` varchar(255) NOT NULL,
  `user_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp(),
  `status` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;

--
-- Dumping data untuk tabel `logs_file_transfer`
--

INSERT INTO `logs_file_transfer` (`id`, `transaction_id`, `user_id`, `product_id`, `timestamp`, `status`) VALUES
(1, '17', 20, 11, '2025-11-20 15:20:35', 1),
(2, '0', 20, 1, '2025-11-20 15:31:51', 1),
(3, '26', 20, 11, '2025-11-20 15:34:47', 1),
(4, '27', 20, 11, '2025-11-20 15:49:43', 1),
(5, '28', 20, 11, '2025-11-20 15:52:50', 1),
(6, '0', 20, 1, '2025-11-20 16:19:25', 1),
(7, '33', 20, 11, '2025-11-21 06:29:34', 1),
(8, '0', 20, 1, '2025-11-21 06:30:20', 1),
(9, '0', 20, 1, '2025-11-21 08:13:27', 1),
(10, '38', 20, 11, '2025-11-21 08:16:39', 1);

-- --------------------------------------------------------

--
-- Struktur dari tabel `logs_request`
--

CREATE TABLE `logs_request` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `identifier_phone` varchar(25) NOT NULL,
  `message` text DEFAULT NULL,
  `last_sent` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `logs_request`
--

INSERT INTO `logs_request` (`id`, `identifier_phone`, `message`, `last_sent`) VALUES
(1, '6285156622644', 'Req OTP: 886154', '2025-11-20 12:38:23'),
(2, '6285156622644', 'Req OTP: 886154', '2025-11-20 12:40:47'),
(3, '6285156622644', 'Create Transaction 1763613890', '2025-11-20 12:45:32'),
(4, '626285156622644', 'Req OTP: 744413', '2025-11-20 12:47:52'),
(5, '6285156622644', 'Req OTP: 386106', '2025-11-20 13:05:52'),
(6, '6285156622644', 'Req OTP: 959778', '2025-11-20 13:06:56'),
(7, '6285156622644', 'Create Transaction 1763644041', '2025-11-20 13:07:23'),
(8, '6285156622644', 'Req OTP: 702489', '2025-11-20 13:28:00'),
(9, '6285156622644', 'Create Transaction 1763645289', '2025-11-20 13:28:12'),
(10, '6285156622644', 'Create Transaction 1763613890', '2025-11-20 13:36:16'),
(11, '6285156622644', 'Req OTP: 271487', '2025-11-20 13:42:06'),
(12, '6285156622644', 'Create Transaction 1763646142', '2025-11-20 13:42:25'),
(13, '6285156622644', 'Req OTP: 948596', '2025-11-20 13:46:40'),
(14, '6285156622644', 'Create Transaction 1763646413', '2025-11-20 13:46:55'),
(15, '6285156622644', 'Req OTP: 847293', '2025-11-20 14:09:19'),
(16, '6285156622644', 'Create Transaction 1763647767', '2025-11-20 14:09:29'),
(17, '6285156622644', 'Create Transaction 1763647767', '2025-11-20 15:22:16'),
(18, '6285156622644', 'Req OTP: 809245', '2025-11-20 15:24:02'),
(19, '6285156622644', 'Create Transaction 1763652251', '2025-11-20 15:24:13'),
(20, '6285156622644', 'Create Transaction 1763652251', '2025-11-20 15:26:14'),
(21, '6285156622644', 'Req OTP: 506412', '2025-11-20 15:27:44'),
(22, '6285156622644', 'Create Transaction 1763652475', '2025-11-20 15:27:57'),
(23, '6285156622644', 'Req OTP: 424625', '2025-11-20 15:29:13'),
(24, '6285156622644', 'Create Transaction 1763652560', '2025-11-20 15:29:23'),
(25, '6285156622644', 'Req OTP: 531449', '2025-11-20 15:31:25'),
(26, '6285156622644', 'Create Transaction 1763652703', '2025-11-20 15:31:45'),
(27, '6285156622644', 'Req OTP: 887610', '2025-11-20 15:32:24'),
(28, '6285156622644', 'Create Transaction 1763652752', '2025-11-20 15:32:34'),
(29, '6285156622644', 'Req OTP: 787229', '2025-11-20 15:33:53'),
(30, '6285156622644', 'Req OTP: 872490', '2025-11-20 15:34:12'),
(31, '6285156622644', 'Create Transaction 1763652859', '2025-11-20 15:34:22'),
(32, '6285156622644', 'Req OTP: 610502', '2025-11-20 15:48:32'),
(33, '6285156622644', 'Req OTP: 632086', '2025-11-20 15:48:53'),
(34, '6285156622644', 'Create Transaction 1763653744', '2025-11-20 15:49:06'),
(35, '6285156622644', 'Req OTP: 344675', '2025-11-20 15:52:17'),
(36, '6285156622644', 'Create Transaction 1763653944', '2025-11-20 15:52:26'),
(37, '6285156622644', 'Req OTP: 744018', '2025-11-20 16:19:05'),
(38, '6285156622644', 'Create Transaction 1763655556', '2025-11-20 16:19:18'),
(39, '6285156622644', 'Req OTP: 453601', '2025-11-21 05:52:09'),
(40, '6285156622644', 'Create Transaction 1763704343', '2025-11-21 05:52:25'),
(41, '6285156622644', 'Req OTP: 727568', '2025-11-21 06:14:37'),
(42, '6285156622644', 'Create Transaction 1763705704', '2025-11-21 06:15:07'),
(43, '6285156622644', 'Req OTP: 929029', '2025-11-21 06:18:05'),
(44, '6285156622644', 'Create Transaction 1763705898', '2025-11-21 06:18:20'),
(45, '6285156622644', 'Req OTP: 425826', '2025-11-21 06:27:43'),
(46, '6285156622644', 'Req OTP: 787291', '2025-11-21 06:28:20'),
(47, '6285156622644', 'Create Transaction 1763706508', '2025-11-21 06:28:30'),
(48, '6285156622644', 'Req OTP: 201020', '2025-11-21 06:30:02'),
(49, '6285156622644', 'Create Transaction 1763706613', '2025-11-21 06:30:15'),
(50, '6285156622644', 'Req OTP: 471802', '2025-11-21 07:31:36'),
(51, '6285156622644', 'Create Transaction 1763710321', '2025-11-21 07:32:03'),
(52, '6285156622644', 'Req OTP: 058710', '2025-11-21 07:37:15'),
(53, '6285156622644', 'Create Transaction 1763710647', '2025-11-21 07:37:29'),
(54, '6285156622644', 'Req OTP: 125251', '2025-11-21 08:02:47'),
(55, '6285156622644', 'Req OTP: 125251', '2025-11-21 08:03:32'),
(56, '6285156622644', 'Req OTP: 706560', '2025-11-21 08:05:37'),
(57, '6285156622644', 'Create Transaction 1763712346', '2025-11-21 08:05:47'),
(58, '6285156622644', 'Req OTP: 914102', '2025-11-21 08:06:56'),
(59, '6285156622644', 'Req OTP: 148309', '2025-11-21 08:13:05'),
(60, '6285156622644', 'Create Transaction 1763712802', '2025-11-21 08:13:22'),
(61, '6285156622644', 'Req OTP: 862108', '2025-11-21 08:15:52'),
(62, '6285156622644', 'Create Transaction 1763712972', '2025-11-21 08:16:12');

-- --------------------------------------------------------

--
-- Struktur dari tabel `product`
--

CREATE TABLE `product` (
  `id` int(11) NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `product_image` text DEFAULT NULL,
  `product_price` decimal(15,2) NOT NULL,
  `product_discount` decimal(15,2) DEFAULT 0.00,
  `product_discount_amount` int(11) DEFAULT 0,
  `product_desc` text DEFAULT NULL,
  `product_viewed` int(11) DEFAULT 0,
  `product_downloaded` int(11) DEFAULT 0,
  `product_create_at` timestamp NULL DEFAULT current_timestamp(),
  `product_category` int(11) DEFAULT NULL,
  `product_path` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;

--
-- Dumping data untuk tabel `product`
--

INSERT INTO `product` (`id`, `product_name`, `product_image`, `product_price`, `product_discount`, `product_discount_amount`, `product_desc`, `product_viewed`, `product_downloaded`, `product_create_at`, `product_category`, `product_path`) VALUES
(1, 'TestB', 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/309809/06/sv01/fnd/IDN/fmt/png/Sepatu-Lari-Pria-Deviate-NITRO™-Elite-3', 150000.00, 0.00, 100, 'Panduan lengkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar bahasa pemrograman Go dari dasar hingga mahir.Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.', 286, 65, '2025-11-20 03:17:23', 3, '/home/user_prinxelio/database/database-product/banner1.png'),
(2, 'Template Website Portofolio', 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/309809/06/sv01/fnd/IDN/fmt/png/Sepatu-Lari-Pria-Deviate-NITRO™-Elite-3', 75000.00, 5000.00, 20, 'Panduan lengkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar bahasa pemrograman Go dari dasar hingga mahir.Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.', 336, 78, '2025-11-20 03:17:23', 1, NULL),
(4, 'Template Website Portofolio', 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/309809/06/sv01/fnd/IDN/fmt/png/Sepatu-Lari-Pria-Deviate-NITRO™-Elite-3', 75000.00, 5000.00, 20, 'Panduan lengkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar bahasa pemrograman Go dari dasar hingga mahir.Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.', 309, 78, '2025-11-20 03:17:23', 1, NULL),
(5, 'Template Website Portofolio', 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/309809/06/sv01/fnd/IDN/fmt/png/Sepatu-Lari-Pria-Deviate-NITRO™-Elite-3', 75000.00, 5000.00, 20, 'Panduan lengkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar bahasa pemrograman Go dari dasar hingga mahir.Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.', 323, 78, '2025-11-20 03:17:23', 1, NULL),
(6, 'Template Website Portofolio', 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/309809/06/sv01/fnd/IDN/fmt/png/Sepatu-Lari-Pria-Deviate-NITRO™-Elite-3', 75000.00, 5000.00, 20, 'Panduan lengkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar bahasa pemrograman Go dari dasar hingga mahir.Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.', 312, 78, '2025-11-20 03:17:23', 1, NULL),
(7, 'Ebook Panduan Go Lengkap', 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/309809/06/sv01/fnd/IDN/fmt/png/Sepatu-Lari-Pria-Deviate-NITRO™-Elite-3', 150000.00, 120000.00, 20, 'Panduan lengkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar bahasa pemrograman Go dari dasar hingga mahir.Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.', 277, 56, '2025-11-20 03:17:23', 3, NULL),
(8, 'Ebook Panduan Go Lengkap', 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/309809/06/sv01/fnd/IDN/fmt/png/Sepatu-Lari-Pria-Deviate-NITRO™-Elite-3', 150000.00, 120000.00, 20, 'Panduan lengkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar bahasa pemrograman Go dari dasar hingga mahir.Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.', 271, 56, '2025-11-20 03:17:23', 3, NULL),
(9, 'Ebook Panduan Go Lengkap', 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/309809/06/sv01/fnd/IDN/fmt/png/Sepatu-Lari-Pria-Deviate-NITRO™-Elite-3', 150000.00, 120000.00, 20, 'Panduan lengkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar bahasa pemrograman Go dari dasar hingga mahir.Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.', 274, 56, '2025-11-20 03:17:23', 3, NULL),
(10, 'Ebook Panduan Go Lengkap', 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/309809/06/sv01/fnd/IDN/fmt/png/Sepatu-Lari-Pria-Deviate-NITRO™-Elite-3', 150000.00, 120000.00, 20, 'Panduan lengkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar bahasa pemrograman Go dari dasar hingga mahir.Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.', 283, 56, '2025-11-20 03:17:23', 3, NULL),
(11, 'TESTA', 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/309809/06/sv01/fnd/IDN/fmt/png/Sepatu-Lari-Pria-Deviate-NITRO™-Elite-3', 75000.00, 5000.00, 20, 'Panduan lengkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar`...ngkap belajar bahasa pemrograman Go dari dasar hingga mahir.Panduan lengkap belajar bahasa pemrograman Go dari dasar hingga mahir.', 335, 84, '2025-11-20 03:17:23', 1, '/home/user_prinxelio/database/database-product/banner1.png');

-- --------------------------------------------------------

--
-- Struktur dari tabel `transactions`
--

CREATE TABLE `transactions` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `base_price` decimal(15,2) NOT NULL,
  `admin_fee` decimal(15,2) DEFAULT 0.00,
  `total_amount` decimal(15,2) NOT NULL,
  `status` enum('UNPAID','PAID','FAILED','EXPIRED','REFUND') NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `merchant_ref` varchar(100) DEFAULT NULL,
  `expired_time` int(11) DEFAULT NULL,
  `reference` varchar(100) NOT NULL,
  `link_qr` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;

--
-- Dumping data untuk tabel `transactions`
--

INSERT INTO `transactions` (`id`, `user_id`, `product_id`, `base_price`, `admin_fee`, `total_amount`, `status`, `created_at`, `updated_at`, `merchant_ref`, `expired_time`, `reference`, `link_qr`) VALUES
(1, 20, 1, 150000.00, 1590.00, 121590.00, 'FAILED', '2025-11-20 05:07:02', '2025-11-20 05:08:16', 'INV-1763615221674', 1763615881, 'DEV-T398713102460ELWW', 'https://tripay.co.id/qr/DEV-T398713102460ELWW'),
(2, 20, 1, 150000.00, 1590.00, 121590.00, 'PAID', '2025-11-20 05:24:10', '2025-11-19 22:24:31', 'INV-1763616249949', 1763616910, 'DEV-T39871310249LEVS9', 'https://tripay.co.id/qr/DEV-T39871310249LEVS9'),
(3, 20, 1, 150000.00, 1590.00, 121590.00, 'PAID', '2025-11-20 05:51:03', '2025-11-19 22:51:38', 'INV-1763617863608', 1763618523, 'DEV-T39871310254HGDG2', 'https://tripay.co.id/qr/DEV-T39871310254HGDG2'),
(4, 20, 2, 75000.00, 1275.00, 76275.00, 'PAID', '2025-11-20 06:06:19', '2025-11-19 23:07:34', 'INV-1763618778807', 1763619438, 'DEV-T39871310257SFARU', 'https://tripay.co.id/qr/DEV-T39871310257SFARU'),
(5, 20, 1, 150000.00, 1590.00, 121590.00, 'PAID', '2025-11-20 07:44:15', '2025-11-20 00:45:32', 'INV-1763624654916', 1763625315, 'DEV-T3987131027723MEZ', 'https://tripay.co.id/qr/DEV-T3987131027723MEZ'),
(6, 20, 1, 150000.00, 1590.00, 121590.00, 'FAILED', '2025-11-20 07:46:12', '2025-11-20 07:46:18', 'INV-1763624772016', 1763625432, 'DEV-T39871310280U7HEB', 'https://tripay.co.id/qr/DEV-T39871310280U7HEB'),
(7, 20, 1, 150000.00, 1590.00, 121590.00, 'FAILED', '2025-11-20 07:46:49', '2025-11-20 07:53:45', 'INV-1763624809222', 1763625469, 'DEV-T39871310281NYHVG', 'https://tripay.co.id/qr/DEV-T39871310281NYHVG'),
(8, 20, 1, 150000.00, 1590.00, 121590.00, 'PAID', '2025-11-20 07:54:03', '2025-11-20 07:54:36', 'INV-1763625243062', 1763625903, 'DEV-T39871310285D565E', 'https://tripay.co.id/qr/DEV-T39871310285D565E'),
(9, 20, 1, 150000.00, 1590.00, 121590.00, 'EXPIRED', '2025-11-20 08:06:10', '2025-11-20 01:07:35', 'INV-1763625970085', 1763626630, 'DEV-T39871310286DNFYY', 'https://tripay.co.id/qr/DEV-T39871310286DNFYY'),
(10, 20, 1, 150000.00, 1590.00, 121590.00, 'PAID', '2025-11-20 08:07:55', '2025-11-20 01:08:16', 'INV-1763626075460', 1763626735, 'DEV-T398713102897ZCCS', 'https://tripay.co.id/qr/DEV-T398713102897ZCCS'),
(11, 20, 1, 150000.00, 1590.00, 121590.00, 'FAILED', '2025-11-20 09:44:53', '2025-11-20 09:50:45', 'INV-1763631893396', 1763632553, 'DEV-T39871310321EC227', 'https://tripay.co.id/qr/DEV-T39871310321EC227'),
(12, 20, 2, 75000.00, 785.00, 5785.00, 'PAID', '2025-11-20 09:52:03', '2025-11-20 09:56:06', 'INV-1763632322879', 1763632983, 'DEV-T39871310322VWXAS', 'https://tripay.co.id/qr/DEV-T39871310322VWXAS'),
(13, 20, 9, 150000.00, 1590.00, 121590.00, 'FAILED', '2025-11-20 10:21:46', '2025-11-20 10:34:13', 'INV-1763634106431', 1763634766, 'DEV-T39871310330IKJDJ', 'https://tripay.co.id/qr/DEV-T39871310330IKJDJ'),
(14, 20, 11, 75000.00, 785.00, 5785.00, 'UNPAID', '2025-11-20 10:36:46', '2025-11-20 10:36:46', 'INV-1763635006225', 1763635666, 'DEV-T39871310336CJPLQ', 'https://tripay.co.id/qr/DEV-T39871310336CJPLQ'),
(15, 20, 2, 75000.00, 785.00, 5785.00, 'EXPIRED', '2025-11-20 10:50:25', '2025-11-20 10:51:47', 'INV-1763635825339', 1763636485, 'DEV-T39871310340WNBI1', 'https://tripay.co.id/qr/DEV-T39871310340WNBI1'),
(16, 20, 10, 150000.00, 1590.00, 121590.00, 'FAILED', '2025-11-20 11:03:35', '2025-11-20 11:03:54', 'INV-1763636615069', 1763637275, 'DEV-T39871310341XBWRF', 'https://tripay.co.id/qr/DEV-T39871310341XBWRF'),
(17, 20, 11, 75000.00, 785.00, 5785.00, 'EXPIRED', '2025-11-20 11:04:24', '2025-11-20 04:06:38', 'INV-1763636664647', 1763637324, 'DEV-T39871310342U2Q46', 'https://tripay.co.id/qr/DEV-T39871310342U2Q46'),
(18, 20, 6, 75000.00, 785.00, 5785.00, 'FAILED', '2025-11-20 13:07:23', '2025-11-20 13:07:43', 'INV-1763644043174', 1763644703, 'DEV-T39871310377ZY8PX', 'https://tripay.co.id/qr/DEV-T39871310377ZY8PX'),
(19, 20, 8, 150000.00, 1590.00, 121590.00, 'PAID', '2025-11-20 13:28:12', '2025-11-20 13:28:48', 'INV-1763645292121', 1763645952, 'DEV-T398713103797G2WV', 'https://tripay.co.id/qr/DEV-T398713103797G2WV'),
(20, 20, 1, 150000.00, 1800.00, 151800.00, 'FAILED', '2025-11-20 13:42:25', '2025-11-20 13:43:38', 'INV-1763646145208', 1763646805, 'DEV-T39871310380K0FV1', 'https://tripay.co.id/qr/DEV-T39871310380K0FV1'),
(21, 20, 2, 75000.00, 785.00, 5785.00, 'FAILED', '2025-11-20 13:46:56', '2025-11-20 13:49:59', 'INV-1763646415671', 1763647075, 'DEV-T398713103819OA9H', 'https://tripay.co.id/qr/DEV-T398713103819OA9H'),
(22, 20, 1, 150000.00, 0.00, 0.00, 'PAID', '2025-11-20 14:09:29', '2025-11-20 14:09:29', '', 0, 'FREE-432dbc1820', ''),
(23, 20, 10, 150000.00, 1590.00, 121590.00, 'PAID', '2025-11-20 15:24:13', '2025-11-20 15:24:27', 'INV-1763652253606', 1763652913, 'DEV-T39871310398PMAOX', 'https://tripay.co.id/qr/DEV-T39871310398PMAOX'),
(24, 20, 1, 150000.00, 0.00, 0.00, 'PAID', '2025-11-20 15:31:52', '2025-11-20 15:31:52', '', 0, 'FREE-6616eccb6b', ''),
(25, 20, 11, 75000.00, 785.00, 5785.00, 'FAILED', '2025-11-20 15:32:34', '2025-11-20 15:33:46', 'INV-1763652754351', 1763653414, 'DEV-T39871310400HUWXA', 'https://tripay.co.id/qr/DEV-T39871310400HUWXA'),
(26, 20, 11, 75000.00, 785.00, 5785.00, 'PAID', '2025-11-20 15:34:22', '2025-11-20 08:34:42', 'INV-1763652862325', 1763653522, 'DEV-T39871310401URMLP', 'https://tripay.co.id/qr/DEV-T39871310401URMLP'),
(27, 20, 11, 75000.00, 785.00, 5785.00, 'PAID', '2025-11-20 15:49:07', '2025-11-20 08:49:43', 'INV-1763653746959', 1763654407, 'DEV-T39871310402U1AAQ', 'https://tripay.co.id/qr/DEV-T39871310402U1AAQ'),
(28, 20, 11, 75000.00, 785.00, 5785.00, 'PAID', '2025-11-20 15:52:27', '2025-11-20 08:52:50', 'INV-1763653947072', 1763654607, 'DEV-T39871310403SI9XN', 'https://tripay.co.id/qr/DEV-T39871310403SI9XN'),
(29, 20, 1, 150000.00, 0.00, 0.00, 'PAID', '2025-11-20 16:19:25', '2025-11-20 16:19:25', '', 0, 'FREE-32250ecbb4', ''),
(30, 20, 11, 75000.00, 785.00, 5785.00, 'FAILED', '2025-11-21 05:52:26', '2025-11-21 05:53:45', 'INV-1763704345956', 1763705006, 'DEV-T39871310603BO8W8', 'https://tripay.co.id/qr/DEV-T39871310603BO8W8'),
(31, 20, 11, 75000.00, 785.00, 5785.00, 'PAID', '2025-11-21 06:15:07', '2025-11-21 06:16:20', 'INV-1763705707230', 1763706367, 'DEV-T39871310608BZHVW', 'https://tripay.co.id/qr/DEV-T39871310608BZHVW'),
(32, 20, 11, 75000.00, 785.00, 5785.00, 'UNPAID', '2025-11-21 06:18:21', '2025-11-21 06:18:21', 'INV-1763705901073', 1763706561, 'DEV-T39871310610QTPRF', 'https://tripay.co.id/qr/DEV-T39871310610QTPRF'),
(33, 20, 11, 75000.00, 785.00, 5785.00, 'PAID', '2025-11-21 06:28:31', '2025-11-20 23:29:35', 'INV-1763706510850', 1763707171, 'DEV-T39871310611G08BD', 'https://tripay.co.id/qr/DEV-T39871310611G08BD'),
(34, 20, 1, 150000.00, 0.00, 0.00, 'PAID', '2025-11-21 06:30:20', '2025-11-21 06:30:20', '', 0, 'FREE-ec9e01491d', ''),
(35, 20, 11, 75000.00, 785.00, 5785.00, 'FAILED', '2025-11-21 07:32:03', '2025-11-21 07:33:25', 'INV-1763710323381', 1763710983, 'DEV-T39871310642DFIIU', 'https://tripay.co.id/qr/DEV-T39871310642DFIIU'),
(36, 20, 11, 75000.00, 785.00, 5785.00, 'REFUND', '2025-11-21 07:37:30', '2025-11-21 00:41:26', 'INV-1763710650112', 1763711310, 'DEV-T39871310644DUSPX', 'https://tripay.co.id/qr/DEV-T39871310644DUSPX'),
(37, 20, 1, 150000.00, 0.00, 0.00, 'PAID', '2025-11-21 08:13:28', '2025-11-21 08:13:28', '', 0, 'FREE-4879661036', ''),
(38, 20, 11, 75000.00, 785.00, 5785.00, 'PAID', '2025-11-21 08:16:13', '2025-11-21 01:16:40', 'INV-1763712972959', 1763713633, 'DEV-T39871310652TBTYM', 'https://tripay.co.id/qr/DEV-T39871310652TBTYM');

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `phone_number` varchar(20) NOT NULL,
  `otp_code` varchar(6) DEFAULT NULL,
  `otp_created_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `last_login` timestamp NULL DEFAULT NULL,
  `status` enum('ALLOW','REJECT') NOT NULL DEFAULT 'ALLOW',
  `banned_time` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;

--
-- Dumping data untuk tabel `users`
--

INSERT INTO `users` (`id`, `phone_number`, `otp_code`, `otp_created_at`, `created_at`, `last_login`, `status`, `banned_time`) VALUES
(1, '6281234567890', '001452', '2025-11-19 21:43:17', '2025-11-20 03:17:57', NULL, 'ALLOW', NULL),
(20, '6285156622644', NULL, '2025-11-21 01:15:51', '2025-11-20 04:11:24', '2025-11-21 08:16:12', 'ALLOW', NULL),
(22, '85156622644', '886154', '2025-11-19 21:12:56', '2025-11-20 04:12:58', NULL, 'ALLOW', NULL),
(26, '6281516622644', '764451', '2025-11-19 21:43:37', '2025-11-20 04:43:39', NULL, 'ALLOW', NULL);

--
-- Indexes for dumped tables
--

--
-- Indeks untuk tabel `category`
--
ALTER TABLE `category`
  ADD PRIMARY KEY (`id`);

--
-- Indeks untuk tabel `logs_file_transfer`
--
ALTER TABLE `logs_file_transfer`
  ADD PRIMARY KEY (`id`);

--
-- Indeks untuk tabel `logs_request`
--
ALTER TABLE `logs_request`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_identifier_phone_logs` (`identifier_phone`);

--
-- Indeks untuk tabel `product`
--
ALTER TABLE `product`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_product_category` (`product_category`);

--
-- Indeks untuk tabel `transactions`
--
ALTER TABLE `transactions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `reference` (`reference`),
  ADD KEY `user_id` (`user_id`),
  ADD KEY `product_id` (`product_id`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `phone_number` (`phone_number`);

--
-- AUTO_INCREMENT untuk tabel yang dibuang
--

--
-- AUTO_INCREMENT untuk tabel `category`
--
ALTER TABLE `category`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT untuk tabel `logs_file_transfer`
--
ALTER TABLE `logs_file_transfer`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

--
-- AUTO_INCREMENT untuk tabel `logs_request`
--
ALTER TABLE `logs_request`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=63;

--
-- AUTO_INCREMENT untuk tabel `product`
--
ALTER TABLE `product`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;

--
-- AUTO_INCREMENT untuk tabel `transactions`
--
ALTER TABLE `transactions`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=39;

--
-- AUTO_INCREMENT untuk tabel `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=129;

--
-- Ketidakleluasaan untuk tabel pelimpahan (Dumped Tables)
--

--
-- Ketidakleluasaan untuk tabel `product`
--
ALTER TABLE `product`
  ADD CONSTRAINT `fk_product_category` FOREIGN KEY (`product_category`) REFERENCES `category` (`id`);

--
-- Ketidakleluasaan untuk tabel `transactions`
--
ALTER TABLE `transactions`
  ADD CONSTRAINT `transactions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  ADD CONSTRAINT `transactions_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
