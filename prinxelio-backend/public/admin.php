<?php
session_start();

// ==================================================================================
// 1. KONFIGURASI & KONEKSI DATABASE
// ==================================================================================

define('DB_HOST', getenv('DB_HOST') ?: '103.87.67.77');
define('DB_PORT', getenv('DB_PORT') ?: '3306');
define('DB_USER', getenv('DB_USER') ?: 'user_prinxelio_digital');
define('DB_PASS', getenv('DB_PASS') ?: '@Digital081297');
define('DB_NAME', getenv('DB_NAME') ?: 'user_prinxelio_digital');
define('ADMIN_USER', getenv('ADMIN_USER') ?: 'admin');
define('ADMIN_PASS', getenv('ADMIN_PASS') ?: 'admin123');

$conn = new mysqli(DB_HOST, DB_USER, DB_PASS, DB_NAME, intval(DB_PORT));
if ($conn->connect_error) {
    die("Koneksi Gagal: " . $conn->connect_error);
}

// Helper Functions
function query($sql) {
    global $conn;
    $result = $conn->query($sql);
    if (!$result) return [];
    $rows = [];
    while ($row = $result->fetch_assoc()) {
        $rows[] = $row;
    }
    return $rows;
}

function escape($data) {
    global $conn;
    return mysqli_real_escape_string($conn, htmlspecialchars($data));
}

function formatRupiah($angka){
    return "Rp " . number_format($angka,0,',','.');
}

function redirect($url) {
    echo "<script>window.location.href='$url';</script>";
    exit;
}

// ==================================================================================
// 2. LOGIC AUTHENTICATION
// ==================================================================================

if (isset($_GET['action']) && $_GET['action'] == 'logout') {
    session_destroy();
    redirect('index.php');
}

if (isset($_POST['login'])) {
    $user = $_POST['username'];
    $pass = $_POST['password'];
    if ($user === ADMIN_USER && $pass === ADMIN_PASS) {
        $_SESSION['admin_logged_in'] = true;
        redirect('index.php');
    } else {
        $error_login = "Username atau Password salah!";
    }
}

// ==================================================================================
// 3. LOGIC POST ACTIONS (CRUD)
// ==================================================================================

// --- USERS ---
if (isset($_POST['toggle_ban_user'])) {
    $id = escape($_POST['user_id']);
    $current_status = escape($_POST['current_status']);
    $new_status = ($current_status == 'ALLOW') ? 'REJECT' : 'ALLOW';
    $conn->query("UPDATE users SET status='$new_status' WHERE id='$id'");
    redirect('?page=users');
}
if (isset($_POST['edit_user'])) {
    $id = escape($_POST['user_id']);
    $phone = escape($_POST['phone_number']);
    $conn->query("UPDATE users SET phone_number='$phone' WHERE id='$id'");
    redirect('?page=users');
}
if (isset($_POST['delete_user'])) {
    $id = escape($_POST['user_id']);
    $conn->query("DELETE FROM users WHERE id='$id'");
    redirect('?page=users');
}

// --- PRODUCTS ---
if (isset($_POST['save_product'])) {
    $name = escape($_POST['product_name']);
    $price = escape($_POST['product_price']);
    $cat = escape($_POST['product_category']);
    $desc = escape($_POST['product_desc']);
    $img = escape($_POST['product_image']);
    $path = escape($_POST['product_path']); // Tambahan Path
    
    if(!empty($_POST['product_id'])) {
        $id = escape($_POST['product_id']);
        $conn->query("UPDATE product SET product_name='$name', product_price='$price', product_category='$cat', product_desc='$desc', product_image='$img', product_path='$path' WHERE id='$id'");
    } else {
        $conn->query("INSERT INTO product (product_name, product_price, product_category, product_desc, product_image, product_path) VALUES ('$name', '$price', '$cat', '$desc', '$img', '$path')");
    }
    redirect('?page=products');
}

if (isset($_POST['delete_product'])) {
    $id = escape($_POST['product_id']);
    $conn->query("DELETE FROM product WHERE id='$id'");
    redirect('?page=products');
}

// --- CATEGORIES ---
if (isset($_POST['save_category'])) {
    $name = escape($_POST['category_name']);
    $color = escape($_POST['category_color']);
    $img = escape($_POST['category_images']); // Sesuai nama kolom DB: category_images
    
    if(!empty($_POST['category_id'])) {
        $id = escape($_POST['category_id']);
        $conn->query("UPDATE category SET category_name='$name', category_color='$color', category_images='$img' WHERE id='$id'");
    } else {
        $conn->query("INSERT INTO category (category_name, category_color, category_images) VALUES ('$name', '$color', '$img')");
    }
    redirect('?page=categories');
}

if (isset($_POST['delete_category'])) {
    $id = escape($_POST['category_id']);
    $conn->query("DELETE FROM category WHERE id='$id'");
    redirect('?page=categories');
}

// --- TRANSACTIONS ---
if (isset($_POST['update_transaction'])) {
    $id = escape($_POST['trx_id']);
    $status = escape($_POST['trx_status']);
    $conn->query("UPDATE transactions SET status='$status' WHERE id='$id'");
    redirect('?page=transactions');
}
if (isset($_POST['delete_transaction'])) {
    $id = escape($_POST['trx_id']);
    $conn->query("DELETE FROM transactions WHERE id='$id'");
    redirect('?page=transactions');
}

// ==================================================================================
// 4. DATA PREPARATION
// ==================================================================================

$page = isset($_GET['page']) ? $_GET['page'] : 'dashboard';

// Dashboard Stats
$total_revenue = query("SELECT SUM(total_amount) as total FROM transactions WHERE status='PAID'")[0]['total'] ?? 0;
$total_users = query("SELECT COUNT(*) as total FROM users")[0]['total'];
$total_products = query("SELECT COUNT(*) as total FROM product")[0]['total'];
$total_transactions = query("SELECT COUNT(*) as total FROM transactions")[0]['total'];

// Chart Data 1: Revenue 7 Days
$chart_revenue_data = query("SELECT DATE(created_at) as date, SUM(total_amount) as total FROM transactions WHERE status='PAID' AND created_at >= DATE(NOW()) - INTERVAL 7 DAY GROUP BY DATE(created_at) ORDER BY date ASC");

// Chart Data 2: Status Distribution
$chart_status_data = query("SELECT status, COUNT(*) as count FROM transactions GROUP BY status");

// Chart Data 3: Top 5 Products Viewed
$chart_top_products = query("SELECT product_name, product_viewed FROM product ORDER BY product_viewed DESC LIMIT 5");

?>
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Prinxelio Admin Panel</title>
    <link rel="icon" type="image/png" href="images/logo.png">
    
    <!-- Fonts & Icons -->
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/remixicon@3.5.0/fonts/remixicon.css" rel="stylesheet">
    
    <!-- Libraries -->
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script src="https://code.jquery.com/jquery-3.7.0.min.js"></script>
    <link rel="stylesheet" href="https://cdn.datatables.net/1.13.6/css/jquery.dataTables.min.css">
    <script src="https://cdn.datatables.net/1.13.6/js/jquery.dataTables.min.js"></script>

    <style>
        /* --- CSS RESET & VARIABLES --- */
        :root {
            --primary: #000000; /* Hitam sesuai request logo */
            --primary-light: #f3f4f6;
            --accent: #487FFF; /* Warna biru untuk tombol/link */
            --bg-body: #F3F4F6;
            --bg-card: #FFFFFF;
            --text-main: #111827;
            --text-muted: #6B7280;
            --border: #E5E7EB;
            --success: #10B981; --success-bg: #D1FAE5;
            --danger: #EF4444; --danger-bg: #FEE2E2;
            --warning: #F59E0B; --warning-bg: #FEF3C7;
            --shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
        }
        
        * { box-sizing: border-box; margin: 0; padding: 0; font-family: "Inter", sans-serif; }
        body { background-color: var(--bg-body); color: var(--text-main); font-size: 0.9rem; }
        a { text-decoration: none; color: inherit; transition: 0.2s; }
        ul { list-style: none; }
        
        /* --- LAYOUT --- */
        .sidebar { position: fixed; top: 0; left: 0; width: 260px; height: 100vh; background: var(--bg-card); border-right: 1px solid var(--border); z-index: 50; }
        .sidebar-header { height: 70px; display: flex; align-items: center; padding: 0 1.5rem; border-bottom: 1px solid var(--border); gap: 12px; }
        .sidebar-header img { height: 32px; width: auto; }
        .sidebar-header span { font-weight: 700; font-size: 1.2rem; color: #000; letter-spacing: -0.5px; }
        
        .sidebar-nav { padding: 1.5rem 1rem; height: calc(100vh - 70px); overflow-y: auto; }
        .nav-link { display: flex; align-items: center; padding: 0.75rem 1rem; color: var(--text-muted); border-radius: 0.5rem; margin-bottom: 0.25rem; font-weight: 500; }
        .nav-link:hover, .nav-link.active { background-color: var(--primary-light); color: var(--primary); }
        .nav-link i { font-size: 1.25rem; margin-right: 0.75rem; }
        
        .main-content { margin-left: 260px; min-height: 100vh; display: flex; flex-direction: column; }
        .topbar { height: 70px; background: var(--bg-card); border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; padding: 0 2rem; position: sticky; top: 0; z-index: 40; }
        .content-body { padding: 2rem; flex-grow: 1; }
        
        /* --- COMPONENTS --- */
        .card { background: var(--bg-card); border-radius: 0.75rem; border: 1px solid var(--border); box-shadow: var(--shadow); margin-bottom: 1.5rem; overflow: hidden; }
        .card-header { padding: 1rem 1.5rem; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; background: #fff; }
        .card-title { font-size: 1rem; font-weight: 600; color: var(--text-main); }
        .card-body { padding: 1.5rem; }
        
        .btn { display: inline-flex; align-items: center; justify-content: center; padding: 0.5rem 1rem; border-radius: 0.375rem; font-weight: 500; border: 1px solid transparent; cursor: pointer; font-size: 0.875rem; gap: 0.5rem; transition: 0.2s; }
        .btn-primary { background: var(--primary); color: #fff; }
        .btn-primary:hover { background: #333; }
        .btn-accent { background: var(--accent); color: #fff; }
        .btn-accent:hover { background: #3b6cd4; }
        .btn-danger { background: var(--danger); color: #fff; }
        .btn-outline { background: transparent; border-color: var(--border); color: var(--text-main); }
        .btn-outline:hover { background: var(--primary-light); }
        .btn-sm { padding: 0.25rem 0.5rem; font-size: 0.75rem; }
        
        .form-control, .form-select { width: 100%; padding: 0.5rem 0.75rem; border: 1px solid var(--border); border-radius: 0.375rem; outline: none; font-size: 0.875rem; margin-bottom: 1rem; }
        .form-control:focus { border-color: var(--accent); ring: 2px solid var(--accent); }
        .form-label { display: block; margin-bottom: 0.375rem; font-weight: 500; font-size: 0.875rem; color: var(--text-main); }
        
        /* --- TABLES --- */
        table.dataTable { border-collapse: collapse !important; width: 100% !important; }
        table.dataTable thead th { background: #F9FAFB; color: var(--text-muted); font-weight: 600; font-size: 0.75rem; text-transform: uppercase; padding: 0.75rem 1rem; border-bottom: 1px solid var(--border); }
        table.dataTable tbody td { padding: 0.75rem 1rem; border-bottom: 1px solid var(--border); color: var(--text-main); vertical-align: middle; font-size: 0.875rem; }
        .dataTables_wrapper .dataTables_filter input { border: 1px solid var(--border); padding: 0.4rem; border-radius: 4px; margin-left: 8px; }
        
        /* --- UTILS --- */
        .d-flex { display: flex; }
        .align-items-center { align-items: center; }
        .justify-between { justify-content: space-between; }
        .gap-2 { gap: 0.5rem; } .gap-4 { gap: 1rem; }
        .grid-4 { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1.5rem; }
        .grid-3 { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1.5rem; }
        .grid-2 { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1.5rem; }
        .w-100 { width: 100%; }
        .text-right { text-align: right; }
        .hidden { display: none !important; }

        .cat-item { position: relative; width: 100%; aspect-ratio: 1 / 1; border-radius: 8px; overflow: hidden; background: var(--primary-light); }
        .cat-thumb { position: absolute; inset: 0; }
        .cat-thumb img { width: 100%; height: 100%; object-fit: cover; display: block; }
        .cat-overlay { position: absolute; inset: 0; background: linear-gradient(to bottom, rgba(255,255,255,0) 0%, currentColor 100%); }
        .cat-name { position: absolute; bottom: 10px; left: 0; right: 0; z-index: 2; font-size: 0.9rem; font-weight: 800; text-transform: uppercase; color: #fff; text-align: center; text-shadow: 0 2px 6px rgba(0,0,0,0.4); }
        
        /* Badges */
        .badge { padding: 2px 8px; border-radius: 9999px; font-size: 0.7rem; font-weight: 600; }
        .badge-success { background: var(--success-bg); color: var(--success); }
        .badge-danger { background: var(--danger-bg); color: var(--danger); }
        .badge-warning { background: var(--warning-bg); color: var(--warning); }
        .badge-gray { background: var(--primary-light); color: var(--text-muted); }

        /* --- PRODUCT & CATEGORY GRID --- */
        .grid-view-container { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 1.5rem; }
        .grid-card { background: #fff; border: 1px solid var(--border); border-radius: 0.5rem; overflow: hidden; transition: transform 0.2s; display: flex; flex-direction: column; }
        .grid-card:hover { transform: translateY(-2px); box-shadow: var(--shadow); }
        .grid-img { height: 140px; width: 100%; object-fit: cover; background-color: #f3f4f6; }
        .grid-body { padding: 1rem; flex-grow: 1; display: flex; flex-direction: column; justify-content: space-between; }
        .grid-title { font-weight: 600; margin-bottom: 0.25rem; font-size: 0.95rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
        .grid-sub { font-size: 0.75rem; color: var(--text-muted); margin-bottom: 0.5rem; }
        .grid-price { font-weight: 700; color: var(--accent); margin-top: auto; }
        
        /* --- MODAL --- */
        .modal { display: none; position: fixed; z-index: 100; left: 0; top: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.5); backdrop-filter: blur(2px); }
        .modal-content { background-color: #fff; margin: 5% auto; padding: 0; border-radius: 0.75rem; width: 90%; max-width: 600px; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1); animation: slideDown 0.3s; }
        @keyframes slideDown { from {transform: translateY(-20px); opacity: 0;} to {transform: translateY(0); opacity: 1;} }
        .modal-body { padding: 1.5rem; max-height: 70vh; overflow-y: auto; }
        .close { float: right; font-size: 1.5rem; cursor: pointer; color: var(--text-muted); }
        
        /* Toggle View Buttons */
        .view-btn { padding: 6px 10px; border: 1px solid var(--border); background: #fff; cursor: pointer; }
        .view-btn.active { background: var(--primary-light); color: var(--primary); border-color: var(--text-muted); }
        .view-btn:first-child { border-top-left-radius: 4px; border-bottom-left-radius: 4px; }
        .view-btn:last-child { border-top-right-radius: 4px; border-bottom-right-radius: 4px; }

        @media (max-width: 768px) {
            .sidebar { transform: translateX(-100%); transition: 0.3s; }
            .sidebar.open { transform: translateX(0); }
            .main-content { margin-left: 0; }
            .grid-4, .grid-3, .grid-2 { grid-template-columns: 1fr; }
        }
    </style>
</head>
<body>

<?php if (!isset($_SESSION['admin_logged_in'])): ?>
    <!-- LOGIN PAGE -->
    <div style="min-height: 100vh; display: flex; align-items: center; justify-content: center; background: var(--bg-body);">
        <div class="card" style="width: 100%; max-width: 400px; padding: 1rem;">
            <div class="card-header" style="justify-content: center; border: none;">
                <div style="text-align: center;">
                    <img src="images/logo.png" alt="Logo" style="height: 50px; margin-bottom: 10px;">
                    <h3 style="color: #000;">Admin Login</h3>
                </div>
            </div>
            <div class="card-body">
                <?php if(isset($error_login)): ?>
                    <div style="background: var(--danger-bg); color: var(--danger); padding: 10px; border-radius: 4px; margin-bottom: 15px; text-align: center;"><?= $error_login ?></div>
                <?php endif; ?>
                <form method="POST">
                    <label class="form-label">Username</label>
                    <input type="text" name="username" class="form-control" required>
                    <label class="form-label">Password</label>
                    <input type="password" name="password" class="form-control" required>
                    <button type="submit" name="login" class="btn btn-primary w-100" style="justify-content: center;">Masuk Dashboard</button>
                </form>
            </div>
        </div>
    </div>
<?php else: ?>
    <!-- ADMIN DASHBOARD LAYOUT -->
    
    <!-- Sidebar -->
    <aside class="sidebar">
        <div class="sidebar-header">
            <img src="images/logo.png" alt="Logo">
            <span>Prinxelio</span>
        </div>
        <nav class="sidebar-nav">
            <a href="?page=dashboard" class="nav-link <?= $page=='dashboard'?'active':'' ?>"><i class="ri-dashboard-line"></i> Dashboard</a>
            <a href="?page=transactions" class="nav-link <?= $page=='transactions'?'active':'' ?>"><i class="ri-file-list-3-line"></i> Transaksi</a>
            <a href="?page=products" class="nav-link <?= $page=='products'?'active':'' ?>"><i class="ri-box-3-line"></i> Produk</a>
            <a href="?page=categories" class="nav-link <?= $page=='categories'?'active':'' ?>"><i class="ri-price-tag-3-line"></i> Kategori</a>
            <a href="?page=users" class="nav-link <?= $page=='users'?'active':'' ?>"><i class="ri-user-line"></i> Users</a>
            <a href="?page=logs" class="nav-link <?= $page=='logs'?'active':'' ?>"><i class="ri-terminal-box-line"></i> System Logs</a>
        </nav>
    </aside>

    <!-- Main Content -->
    <main class="main-content">
        <header class="topbar">
            <h2 style="font-size: 1.1rem; font-weight: 600;">Admin Panel</h2>
            <a href="?action=logout" class="btn btn-sm btn-outline" style="color: var(--danger); border-color: var(--danger);"><i class="ri-logout-box-r-line"></i> Logout</a>
        </header>

        <div class="content-body">
            
            <?php if ($page == 'dashboard'): ?>
                <!-- DASHBOARD -->
                <div class="grid-4 mb-4" style="margin-bottom: 2rem;">
                    <div class="card" style="margin:0;">
                        <div class="card-body">
                            <span class="text-muted" style="font-size: 0.8rem;">Total Pendapatan</span>
                            <h3 style="font-size: 1.5rem; margin-top: 5px;"><?= formatRupiah($total_revenue) ?></h3>
                            <div style="margin-top: 10px; font-size: 0.8rem; color: var(--success);"><i class="ri-arrow-up-line"></i> Paid Transactions</div>
                        </div>
                    </div>
                    <div class="card" style="margin:0;">
                        <div class="card-body">
                            <span class="text-muted" style="font-size: 0.8rem;">Total Transaksi</span>
                            <h3 style="font-size: 1.5rem; margin-top: 5px;"><?= $total_transactions ?></h3>
                        </div>
                    </div>
                    <div class="card" style="margin:0;">
                        <div class="card-body">
                            <span class="text-muted" style="font-size: 0.8rem;">Total Produk</span>
                            <h3 style="font-size: 1.5rem; margin-top: 5px;"><?= $total_products ?></h3>
                        </div>
                    </div>
                    <div class="card" style="margin:0;">
                        <div class="card-body">
                            <span class="text-muted" style="font-size: 0.8rem;">Total User</span>
                            <h3 style="font-size: 1.5rem; margin-top: 5px;"><?= $total_users ?></h3>
                        </div>
                    </div>
                </div>

                <div class="grid-3" style="margin-bottom: 2rem;">
                    <!-- Chart Revenue -->
                    <div class="card" style="grid-column: span 2; margin:0;">
                        <div class="card-header"><h6 class="card-title">Pendapatan 7 Hari Terakhir</h6></div>
                        <div class="card-body"><canvas id="revenueChart" style="max-height: 300px;"></canvas></div>
                    </div>
                    <!-- Chart Status -->
                    <div class="card" style="margin:0;">
                        <div class="card-header"><h6 class="card-title">Status Transaksi</h6></div>
                        <div class="card-body"><canvas id="statusChart" style="max-height: 200px;"></canvas></div>
                    </div>
                </div>

                <div class="grid-2">
                    <!-- Chart Top Products -->
                    <div class="card" style="margin:0;">
                        <div class="card-header"><h6 class="card-title">5 Produk Terpopuler (Dilihat)</h6></div>
                        <div class="card-body"><canvas id="topProductChart" style="max-height: 250px;"></canvas></div>
                    </div>
                    <!-- Recent Transactions Table -->
                    <div class="card" style="margin:0;">
                        <div class="card-header"><h6 class="card-title">Transaksi Terakhir</h6></div>
                        <div class="card-body" style="padding: 0;">
                            <table class="table" style="width: 100%; font-size: 0.85rem;">
                                <thead style="background: #f9fafb;">
                                    <tr><th style="padding: 10px 15px; text-align: left;">Ref</th><th style="padding: 10px 15px; text-align: left;">Total</th><th style="padding: 10px 15px; text-align: left;">Status</th></tr>
                                </thead>
                                <tbody>
                                    <?php 
                                    $recent_trx = query("SELECT reference, total_amount, status FROM transactions ORDER BY id DESC LIMIT 5");
                                    foreach($recent_trx as $rt):
                                        $badge = match($rt['status']) { 'PAID'=>'badge-success', 'FAILED'=>'badge-danger', 'UNPAID'=>'badge-warning', default=>'badge-gray' };
                                    ?>
                                    <tr>
                                        <td style="padding: 10px 15px; border-bottom: 1px solid var(--border);"><?= $rt['reference'] ?></td>
                                        <td style="padding: 10px 15px; border-bottom: 1px solid var(--border);"><?= formatRupiah($rt['total_amount']) ?></td>
                                        <td style="padding: 10px 15px; border-bottom: 1px solid var(--border);"><span class="badge <?= $badge ?>"><?= $rt['status'] ?></span></td>
                                    </tr>
                                    <?php endforeach; ?>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>

                <script>
                    // Config Chart Revenue
                    new Chart(document.getElementById('revenueChart'), {
                        type: 'line',
                        data: {
                            labels: <?= json_encode(array_column($chart_revenue_data, 'date')) ?>,
                            datasets: [{ label: 'Pendapatan', data: <?= json_encode(array_column($chart_revenue_data, 'total')) ?>, borderColor: '#487FFF', backgroundColor: 'rgba(72, 127, 255, 0.1)', tension: 0.3, fill: true }]
                        }, options: { responsive: true, plugins: { legend: { display: false } } }
                    });
                    // Config Chart Status
                    new Chart(document.getElementById('statusChart'), {
                        type: 'doughnut',
                        data: {
                            labels: <?= json_encode(array_column($chart_status_data, 'status')) ?>,
                            datasets: [{ data: <?= json_encode(array_column($chart_status_data, 'count')) ?>, backgroundColor: ['#10B981', '#EF4444', '#F59E0B', '#9CA3AF', '#487FFF'], borderWidth: 0 }]
                        }, options: { responsive: true, plugins: { legend: { position: 'bottom' } } }
                    });
                    // Config Chart Top Products
                    new Chart(document.getElementById('topProductChart'), {
                        type: 'bar',
                        data: {
                            labels: <?= json_encode(array_column($chart_top_products, 'product_name')) ?>,
                            datasets: [{ label: 'Views', data: <?= json_encode(array_column($chart_top_products, 'product_viewed')) ?>, backgroundColor: '#487FFF', borderRadius: 4 }]
                        }, options: { indexAxis: 'y', responsive: true, plugins: { legend: { display: false } } }
                    });
                </script>

            <?php elseif ($page == 'transactions'): ?>
                <!-- TRANSACTIONS PAGE -->
                <div class="card">
                    <div class="card-header"><h6 class="card-title">Data Transaksi</h6></div>
                    <div class="card-body">
                        <table class="table datatable">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Reference</th>
                                    <th>User</th>
                                    <th>Produk</th>
                                    <th>Harga</th>
                                    <th>Status</th>
                                    <th>Tanggal</th>
                                    <th>Aksi</th>
                                </tr>
                            </thead>
                            <tbody>
                                <?php 
                                $all_trx = query("SELECT t.*, u.phone_number, p.product_name FROM transactions t LEFT JOIN users u ON t.user_id = u.id LEFT JOIN product p ON t.product_id = p.id ORDER BY t.id DESC");
                                foreach($all_trx as $trx): 
                                    $statusClass = match($trx['status']) { 'PAID' => 'badge-success', 'FAILED' => 'badge-danger', 'UNPAID' => 'badge-warning', default => 'badge-gray' };
                                ?>
                                <tr>
                                    <td><?= $trx['id'] ?></td>
                                    <td>
                                        <div style="font-weight: 600;"><?= $trx['reference'] ?></div>
                                        <small class="text-muted" style="font-size: 0.7rem;"><?= $trx['merchant_ref'] ?></small>
                                    </td>
                                    <td><?= $trx['phone_number'] ?></td>
                                    <td><?= $trx['product_name'] ?></td>
                                    <td><?= formatRupiah($trx['total_amount']) ?></td>
                                    <td><span class="badge <?= $statusClass ?>"><?= $trx['status'] ?></span></td>
                                    <td><?= $trx['created_at'] ?></td>
                                    <td>
                                        <div class="d-flex gap-2">
                                            <button class="btn btn-sm btn-outline" onclick="openTrxModal('<?= $trx['id'] ?>', '<?= $trx['status'] ?>')"><i class="ri-edit-line"></i></button>
                                            <form method="POST" onsubmit="return confirm('Hapus transaksi ini?');">
                                                <input type="hidden" name="trx_id" value="<?= $trx['id'] ?>">
                                                <button type="submit" name="delete_transaction" class="btn btn-sm btn-outline" style="color: var(--danger); border-color: var(--danger);"><i class="ri-delete-bin-line"></i></button>
                                            </form>
                                        </div>
                                    </td>
                                </tr>
                                <?php endforeach; ?>
                            </tbody>
                        </table>
                    </div>
                </div>

                <!-- Modal Edit Transaksi -->
                <div id="trxModal" class="modal">
                    <div class="modal-content">
                        <div class="card-header"><h5 class="card-title">Edit Status Transaksi</h5><span class="close" onclick="closeModal('trxModal')">&times;</span></div>
                        <div class="modal-body">
                            <form method="POST">
                                <input type="hidden" name="trx_id" id="edit_trx_id">
                                <label class="form-label">Status</label>
                                <select name="trx_status" id="edit_trx_status" class="form-select">
                                    <option value="UNPAID">UNPAID</option>
                                    <option value="PAID">PAID</option>
                                    <option value="FAILED">FAILED</option>
                                    <option value="EXPIRED">EXPIRED</option>
                                    <option value="REFUND">REFUND</option>
                                </select>
                                <button type="submit" name="update_transaction" class="btn btn-primary w-100">Update Status</button>
                            </form>
                        </div>
                    </div>
                </div>

            <?php elseif ($page == 'products'): ?>
                <!-- PRODUCTS PAGE -->
                <div class="d-flex justify-between align-items-center mb-4" style="margin-bottom: 1.5rem;">
                    <div class="d-flex">
                        <button class="view-btn active" id="btn-prod-grid" onclick="toggleView('product', 'grid')"><i class="ri-grid-fill"></i> Grid</button>
                        <button class="view-btn" id="btn-prod-table" onclick="toggleView('product', 'table')"><i class="ri-table-line"></i> Tabel</button>
                    </div>
                    <button class="btn btn-accent" onclick="openProductModal()"><i class="ri-add-line"></i> Tambah Produk</button>
                </div>

                <?php $products = query("SELECT p.*, c.category_name FROM product p LEFT JOIN category c ON p.product_category = c.id ORDER BY p.id DESC"); ?>

                <!-- VIEW: GRID -->
                <div id="product-grid-view" class="grid-view-container">
                    <?php foreach($products as $prod): ?>
                    <div class="grid-card">
                        <img src="<?= $prod['product_image'] ?>" class="grid-img" alt="img" onerror="this.src='https://via.placeholder.com/300?text=No+Image'">
                        <div class="grid-body">
                            <div>
                                <div class="badge badge-gray" style="margin-bottom: 5px; display: inline-block;"><?= $prod['category_name'] ?></div>
                                <h4 class="grid-title"><?= $prod['product_name'] ?></h4>
                                <div class="grid-sub"><i class="ri-eye-line"></i> <?= $prod['product_viewed'] ?> views</div>
                            </div>
                            <div class="d-flex justify-between align-items-center mt-2">
                                <div class="grid-price"><?= formatRupiah($prod['product_price']) ?></div>
                                <div class="d-flex gap-2">
                                    <button class="btn btn-sm btn-outline" onclick='openProductModal(<?= json_encode($prod) ?>)'><i class="ri-pencil-line"></i></button>
                                    <form method="POST" onsubmit="return confirm('Hapus?');">
                                        <input type="hidden" name="product_id" value="<?= $prod['id'] ?>">
                                        <button type="submit" name="delete_product" class="btn btn-sm btn-outline" style="color: var(--danger); border-color: var(--danger);"><i class="ri-delete-bin-line"></i></button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </div>
                    <?php endforeach; ?>
                </div>

                <!-- VIEW: TABLE (No Image, With Path) -->
                <div id="product-table-view" class="card hidden">
                    <div class="card-body">
                        <table class="table datatable">
                            <thead><tr><th>ID</th><th>Nama Produk</th><th>Kategori</th><th>Harga</th><th>Path File</th><th>Stats</th><th>Aksi</th></tr></thead>
                            <tbody>
                                <?php foreach($products as $prod): ?>
                                <tr>
                                    <td><?= $prod['id'] ?></td>
                                    <td><?= $prod['product_name'] ?></td>
                                    <td><?= $prod['category_name'] ?></td>
                                    <td><?= formatRupiah($prod['product_price']) ?></td>
                                    <td><code style="font-size: 0.75rem; background: #eee; padding: 2px 4px; border-radius: 3px;"><?= $prod['product_path'] ?></code></td>
                                    <td><small><?= $prod['product_viewed'] ?> views / <?= $prod['product_downloaded'] ?> DL</small></td>
                                    <td>
                                        <div class="d-flex gap-2">
                                            <button class="btn btn-sm btn-outline" onclick='openProductModal(<?= json_encode($prod) ?>)'><i class="ri-pencil-line"></i></button>
                                            <form method="POST" onsubmit="return confirm('Hapus?');">
                                                <input type="hidden" name="product_id" value="<?= $prod['id'] ?>">
                                                <button type="submit" name="delete_product" class="btn btn-sm btn-outline" style="color: var(--danger);"><i class="ri-delete-bin-line"></i></button>
                                            </form>
                                        </div>
                                    </td>
                                </tr>
                                <?php endforeach; ?>
                            </tbody>
                        </table>
                    </div>
                </div>

                <!-- Modal Product -->
                <div id="productModal" class="modal">
                    <div class="modal-content">
                        <div class="card-header"><h5 class="card-title" id="prodModalTitle">Tambah Produk</h5><span class="close" onclick="closeModal('productModal')">&times;</span></div>
                        <div class="modal-body">
                            <form method="POST">
                                <input type="hidden" name="product_id" id="prod_id">
                                <div class="grid-2" style="gap: 1rem; margin-bottom: 0;">
                                    <div><label class="form-label">Nama Produk</label><input type="text" name="product_name" id="prod_name" class="form-control" required></div>
                                    <div><label class="form-label">Harga</label><input type="number" name="product_price" id="prod_price" class="form-control" required></div>
                                </div>
                                <label class="form-label">Kategori</label>
                                <select name="product_category" id="prod_cat" class="form-select">
                                    <?php foreach(query("SELECT * FROM category") as $c) echo "<option value='{$c['id']}'>{$c['category_name']}</option>"; ?>
                                </select>
                                <label class="form-label">URL Gambar (Banner)</label>
                                <input type="text" name="product_image" id="prod_img" class="form-control">
                                
                                <label class="form-label">Path File (Lokasi File Produk)</label>
                                <input type="text" name="product_path" id="prod_path" class="form-control" placeholder="/home/user/database/...">
                                
                                <label class="form-label">Deskripsi</label>
                                <textarea name="product_desc" id="prod_desc" class="form-control" rows="3"></textarea>
                                
                                <button type="submit" name="save_product" class="btn btn-primary w-100">Simpan Produk</button>
                            </form>
                        </div>
                    </div>
                </div>

            <?php elseif ($page == 'categories'): ?>
                <!-- CATEGORIES PAGE -->
                <div class="d-flex justify-between align-items-center mb-4" style="margin-bottom: 1.5rem;">
                    <div class="d-flex">
                        <button class="view-btn active" id="btn-cat-grid" onclick="toggleView('cat', 'grid')"><i class="ri-grid-fill"></i> Grid</button>
                        <button class="view-btn" id="btn-cat-table" onclick="toggleView('cat', 'table')"><i class="ri-table-line"></i> Tabel</button>
                    </div>
                    <button class="btn btn-accent" onclick="openCatModal()"><i class="ri-add-line"></i> Tambah Kategori</button>
                </div>

                <?php $categories = query("SELECT * FROM category"); ?>

                <!-- VIEW: GRID (With Image) -->
                <div id="cat-grid-view" class="grid-view-container">
                    <?php foreach($categories as $cat): ?>
                    <div class="grid-card">
                        <div class="cat-item">
                            <div class="cat-thumb">
                                <img src="<?= (str_starts_with($cat['category_images'], 'http://') || str_starts_with($cat['category_images'], 'https://')) 
                                    ? $cat['category_images'] 
                                    : (str_starts_with(ltrim($cat['category_images'], '/'), 'images/') 
                                        ? ltrim($cat['category_images'], '/') 
                                        : ('images/category/' . ltrim($cat['category_images'], '/'))) 
                                ?>" alt="<?= $cat['category_name'] ?>" onerror="this.src='https://placehold.co/300'">
                            </div>
                            <div class="cat-overlay" style="color: <?= $cat['category_color'] ?: '#000' ?>"></div>
                            <div class="cat-name"><?= $cat['category_name'] ?></div>
                        </div>
                        <div class="grid-body" style="flex-grow: 0;">
                            <div class="d-flex justify-between align-items-center">
                                <h4 class="grid-title"><?= $cat['category_name'] ?></h4>
                                <div class="d-flex gap-2">
                                    <button class="btn btn-sm btn-outline" onclick='openCatModal(<?= json_encode($cat) ?>)'><i class="ri-pencil-line"></i></button>
                                    <form method="POST" onsubmit="return confirm('Hapus?');">
                                        <input type="hidden" name="category_id" value="<?= $cat['id'] ?>">
                                        <button type="submit" name="delete_category" class="btn btn-sm btn-outline" style="color: var(--danger);"><i class="ri-delete-bin-line"></i></button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </div>
                    <?php endforeach; ?>
                </div>

                <!-- VIEW: TABLE (Path Text, No Image Render) -->
                <div id="cat-table-view" class="card hidden">
                    <div class="card-body">
                        <table class="table datatable">
                            <thead><tr><th>ID</th><th>Nama</th><th>Warna</th><th>Image Path</th><th>Dibuat</th><th>Aksi</th></tr></thead>
                            <tbody>
                                <?php foreach($categories as $cat): ?>
                                <tr>
                                    <td><?= $cat['id'] ?></td>
                                    <td><?= $cat['category_name'] ?></td>
                                    <td><span style="display:inline-block; width:15px; height:15px; background:<?= $cat['category_color'] ?>; border-radius:50%; vertical-align:middle;"></span> <?= $cat['category_color'] ?></td>
                                    <td><code style="font-size: 0.75rem; background: #eee; padding: 2px 4px; border-radius: 3px;"><?= $cat['category_images'] ?></code></td>
                                    <td><?= $cat['category_create_at'] ?></td>
                                    <td>
                                        <div class="d-flex gap-2">
                                            <button class="btn btn-sm btn-outline" onclick='openCatModal(<?= json_encode($cat) ?>)'><i class="ri-pencil-line"></i></button>
                                            <form method="POST" onsubmit="return confirm('Hapus?');">
                                                <input type="hidden" name="category_id" value="<?= $cat['id'] ?>">
                                                <button type="submit" name="delete_category" class="btn btn-sm btn-outline" style="color: var(--danger);"><i class="ri-delete-bin-line"></i></button>
                                            </form>
                                        </div>
                                    </td>
                                </tr>
                                <?php endforeach; ?>
                            </tbody>
                        </table>
                    </div>
                </div>

                <!-- Modal Category -->
                <div id="catModal" class="modal">
                    <div class="modal-content">
                        <div class="card-header"><h5 class="card-title" id="catModalTitle">Tambah Kategori</h5><span class="close" onclick="closeModal('catModal')">&times;</span></div>
                        <div class="modal-body">
                            <form method="POST">
                                <input type="hidden" name="category_id" id="cat_id">
                                <label class="form-label">Nama Kategori</label>
                                <input type="text" name="category_name" id="cat_name" class="form-control" required>
                                <label class="form-label">Warna (Hex)</label>
                                <input type="color" name="category_color" id="cat_color" class="form-control" style="height: 40px; padding: 2px;" value="#000000">
                                <label class="form-label">URL Gambar / Path</label>
                                <input type="text" name="category_images" id="cat_img" class="form-control" placeholder="/images/category/...">
                                <button type="submit" name="save_category" class="btn btn-primary w-100">Simpan</button>
                            </form>
                        </div>
                    </div>
                </div>

            <?php elseif ($page == 'users'): ?>
                <!-- USERS PAGE -->
                <div class="card">
                    <div class="card-header"><h6 class="card-title">Manajemen User</h6></div>
                    <div class="card-body">
                        <table class="table datatable">
                            <thead><tr><th>ID</th><th>No. HP</th><th>Status</th><th>Terdaftar</th><th>Aksi</th></tr></thead>
                            <tbody>
                                <?php foreach(query("SELECT * FROM users ORDER BY id DESC") as $u): ?>
                                <tr>
                                    <td><?= $u['id'] ?></td>
                                    <td><?= $u['phone_number'] ?></td>
                                    <td><span class="badge <?= $u['status'] == 'ALLOW' ? 'badge-success' : 'badge-danger' ?>"><?= $u['status'] ?></span></td>
                                    <td><?= $u['created_at'] ?></td>
                                    <td>
                                        <div class="d-flex gap-2">
                                            <button class="btn btn-sm btn-outline" onclick="openUserModal('<?= $u['id'] ?>', '<?= $u['phone_number'] ?>')"><i class="ri-pencil-line"></i></button>
                                            <form method="POST">
                                                <input type="hidden" name="user_id" value="<?= $u['id'] ?>">
                                                <input type="hidden" name="current_status" value="<?= $u['status'] ?>">
                                                <button type="submit" name="toggle_ban_user" class="btn btn-sm <?= $u['status'] == 'ALLOW' ? 'btn-outline' : 'btn-primary' ?>" style="<?= $u['status'] == 'ALLOW' ? 'color:var(--danger); border-color:var(--danger);' : '' ?>">
                                                    <?= $u['status'] == 'ALLOW' ? 'Ban' : 'Unban' ?>
                                                </button>
                                            </form>
                                            <form method="POST" onsubmit="return confirm('Hapus user permanen?');">
                                                <input type="hidden" name="user_id" value="<?= $u['id'] ?>">
                                                <button type="submit" name="delete_user" class="btn btn-sm btn-danger"><i class="ri-delete-bin-line"></i></button>
                                            </form>
                                        </div>
                                    </td>
                                </tr>
                                <?php endforeach; ?>
                            </tbody>
                        </table>
                    </div>
                </div>

                <!-- Modal User -->
                <div id="userModal" class="modal">
                    <div class="modal-content">
                        <div class="card-header"><h5 class="card-title">Edit User</h5><span class="close" onclick="closeModal('userModal')">&times;</span></div>
                        <div class="modal-body">
                            <form method="POST">
                                <input type="hidden" name="user_id" id="edit_user_id">
                                <label class="form-label">No. HP</label>
                                <input type="text" name="phone_number" id="edit_user_phone" class="form-control" required>
                                <button type="submit" name="edit_user" class="btn btn-primary w-100">Update</button>
                            </form>
                        </div>
                    </div>
                </div>

            <?php elseif ($page == 'logs'): ?>
                <!-- LOGS PAGE -->
                <div class="grid-2">
                    <div class="card">
                        <div class="card-header"><h6 class="card-title">Request Logs (OTP)</h6></div>
                        <div class="card-body">
                            <table class="table datatable">
                                <thead><tr><th>Identifier</th><th>Message</th></tr></thead>
                                <tbody>
                                    <?php foreach(query("SELECT * FROM logs_request ORDER BY id DESC LIMIT 100") as $lr): ?>
                                    <tr>
                                        <td><div style="font-weight: 600;"><?= $lr['identifier_phone'] ?></div><small class="text-muted"><?= $lr['last_sent'] ?></small></td>
                                        <td><?= $lr['message'] ?></td>
                                    </tr>
                                    <?php endforeach; ?>
                                </tbody>
                            </table>
                        </div>
                    </div>
                    <div class="card">
                        <div class="card-header"><h6 class="card-title">File Download Logs</h6></div>
                        <div class="card-body">
                            <table class="table datatable">
                                <thead><tr><th>User</th><th>Produk</th><th>Waktu</th></tr></thead>
                                <tbody>
                                    <?php foreach(query("SELECT l.*, u.phone_number, p.product_name FROM logs_file_transfer l LEFT JOIN users u ON l.user_id = u.id LEFT JOIN product p ON l.product_id = p.id ORDER BY l.id DESC LIMIT 100") as $lf): ?>
                                    <tr>
                                        <td><?= $lf['phone_number'] ?></td>
                                        <td><?= $lf['product_name'] ?></td>
                                        <td><?= $lf['timestamp'] ?></td>
                                    </tr>
                                    <?php endforeach; ?>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            <?php endif; ?>

        </div>
    </main>

    <!-- JAVASCRIPT LOGIC -->
    <script>
        // Initialize DataTables
        $(document).ready(function() {
            $('.datatable').DataTable({
                "pageLength": 10,
                "lengthChange": false,
                "language": { "search": "", "searchPlaceholder": "Cari data...", "paginate": { "next": ">", "previous": "<" }, "info": "Menampilkan _START_ - _END_ dari _TOTAL_" }
            });
        });

        // Modal Functions
        function closeModal(id) { document.getElementById(id).style.display = "none"; }
        window.onclick = function(event) { if (event.target.classList.contains('modal')) { event.target.style.display = "none"; } }

        function openTrxModal(id, status) {
            document.getElementById('edit_trx_id').value = id;
            document.getElementById('edit_trx_status').value = status;
            document.getElementById('trxModal').style.display = "block";
        }

        function openProductModal(data = null) {
            if(data) {
                document.getElementById('prodModalTitle').innerText = 'Edit Produk';
                document.getElementById('prod_id').value = data.id;
                document.getElementById('prod_name').value = data.product_name;
                document.getElementById('prod_price').value = data.product_price;
                document.getElementById('prod_cat').value = data.product_category;
                document.getElementById('prod_img').value = data.product_image;
                document.getElementById('prod_path').value = data.product_path; // Load Path
                document.getElementById('prod_desc').value = data.product_desc;
            } else {
                document.getElementById('prodModalTitle').innerText = 'Tambah Produk';
                document.getElementById('prod_id').value = '';
                document.getElementById('prod_name').value = '';
                document.getElementById('prod_price').value = '';
                document.getElementById('prod_img').value = '';
                document.getElementById('prod_path').value = '';
                document.getElementById('prod_desc').value = '';
            }
            document.getElementById('productModal').style.display = "block";
        }

        function openCatModal(data = null) {
            if(data) {
                document.getElementById('catModalTitle').innerText = 'Edit Kategori';
                document.getElementById('cat_id').value = data.id;
                document.getElementById('cat_name').value = data.category_name;
                document.getElementById('cat_color').value = data.category_color;
                document.getElementById('cat_img').value = data.category_images || ''; // Fix column name
            } else {
                document.getElementById('catModalTitle').innerText = 'Tambah Kategori';
                document.getElementById('cat_id').value = '';
                document.getElementById('cat_name').value = '';
                document.getElementById('cat_color').value = '#000000';
                document.getElementById('cat_img').value = '';
            }
            document.getElementById('catModal').style.display = "block";
        }

        function openUserModal(id, phone) {
            document.getElementById('edit_user_id').value = id;
            document.getElementById('edit_user_phone').value = phone;
            document.getElementById('userModal').style.display = "block";
        }

        // Toggle Grid/Table View
        function toggleView(type, view) {
            // type: 'product' or 'cat'
            // view: 'grid' or 'table'
            
            const gridContainer = document.getElementById(type + '-grid-view');
            const tableContainer = document.getElementById(type + '-table-view');
            const btnGrid = document.getElementById('btn-' + type + '-grid');
            const btnTable = document.getElementById('btn-' + type + '-table');

            if(view === 'grid') {
                gridContainer.classList.remove('hidden');
                tableContainer.classList.add('hidden');
                btnGrid.classList.add('active');
                btnTable.classList.remove('active');
            } else {
                gridContainer.classList.add('hidden');
                tableContainer.classList.remove('hidden');
                btnGrid.classList.remove('active');
                btnTable.classList.add('active');
            }
        }
    </script>
<?php endif; ?>

</body>
</html>
