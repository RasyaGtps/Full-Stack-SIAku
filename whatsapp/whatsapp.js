const { Client, LocalAuth } = require('whatsapp-web.js');
const qrcode = require('qrcode-terminal');
const fs = require('fs');
const Jimp = require('jimp');
const axios = require('axios');

let client;
let qrCode = null;
let isReady = false;
let botInfo = {
    number: null,
    name: null
};
const owners = new Set();
const blockedUsers = new Set();
const pendingVerification = new Map();
const userSessions = new Map(); // Store logged in users: phoneNumber -> {nim, nama, token, role}

async function initWhatsApp() {
    console.log('🔧 Initializing WhatsApp Client...\n');

    client = new Client({
        authStrategy: new LocalAuth({
            dataPath: './auth_session'
        }),
        puppeteer: {
            headless: true,
            args: [
                '--no-sandbox',
                '--disable-setuid-sandbox',
                '--disable-dev-shm-usage',
                '--disable-accelerated-2d-canvas',
                '--no-first-run',
                '--no-zygote',
                '--disable-gpu'
            ]
        }
    });

    client.on('qr', (qr) => {
        qrCode = qr;
        console.log('\n╔═══════════════════════════════════════════════════════╗');
        console.log('║            📱 SCAN QR CODE DENGAN WHATSAPP           ║');
        console.log('╚═══════════════════════════════════════════════════════╝\n');
        qrcode.generate(qr, { small: true });
        console.log('\n✅ QR Code berhasil digenerate!');
        console.log('🔗 Atau akses via API: http://localhost:3000/api/wa/qr-code');
        console.log('⏰ Scan sekarang sebelum expired!\n');
    });

    client.on('ready', () => {
        isReady = true;
        qrCode = null;
        botInfo.number = client.info.wid.user;
        botInfo.name = client.info.pushname;
        console.log('\n✅ WhatsApp Client Ready!');
        console.log(`📱 Number: ${botInfo.number}`);
        console.log(`📝 Name: ${botInfo.name}`);
        console.log('🎉 Bot siap digunakan!\n');
        loadData();
    });
        
    client.on('authenticated', () => {
        console.log('✅ Authenticated!');
    });

    // Authentication failure
    client.on('auth_failure', (msg) => {
        console.error('❌ Authentication failure:', msg);
        isReady = false;
    });

    // Disconnected event
    client.on('disconnected', (reason) => {
        console.log('⚠️  Disconnected:', reason);
        isReady = false;
        qrCode = null;
    });

    // Handle incoming messages
    client.on('message', async (msg) => {
        try {
            if (msg.fromMe || msg.from.includes('@g.us')) return; // Skip own messages and groups

            const phoneNumber = msg.from.replace('@c.us', '');
            const messageText = msg.body || '';

            console.log(`📩 [${phoneNumber}]: ${messageText}`);

            // Check if blocked
            if (blockedUsers.has(phoneNumber)) {
                await msg.reply('❌ Anda telah di-block oleh owner bot.');
                return;
            }

            // Handle image messages for profile picture
            if (msg.hasMedia && msg.type === 'image') {
                const caption = msg.body || '';
                if ((caption === 'setpp' || caption === '/setpp') && isOwner(phoneNumber)) {
                    await handleSetProfilePicture(msg);
                    return;
                }
            }

            // Handle text commands
            if (messageText) {
                await handleCommand(msg, phoneNumber, messageText);
            }
        } catch (error) {
            console.error('Error handling message:', error);
        }
    });

    // Initialize client
    console.log('🚀 Starting WhatsApp client...\n');
    await client.initialize();
}

// Handle commands
async function handleCommand(msg, phoneNumber, message) {
    const args = message.trim().split(' ');
    const command = args[0].toLowerCase();

    switch (command) {
        case '/jadiowner':
            await handleJadiOwner(msg, phoneNumber);
            break;
        case '/cekowner':
            await handleCekOwner(msg, phoneNumber);
            break;
        case '/keluarowner':
            await handleKeluarOwner(msg, phoneNumber);
            break;
        case '/menu':
        case '/help':
            await handleMenu(msg, phoneNumber);
            break;
        case '/gantinama':
            if (isOwner(phoneNumber)) {
                const newName = args.slice(1).join(' ');
                if (newName) await handleGantiNama(msg, newName);
                else await msg.reply('❌ Format: /gantinama [nama baru]\n\nContoh: /gantinama SIAku Bot');
            }
            break;
        case '/gantipp':
            if (isOwner(phoneNumber)) {
                await msg.reply('📷 *CARA GANTI PROFILE PICTURE*\n\n' +
                    'Step by step:\n\n' +
                    '1️⃣ Pilih foto dari galeri\n' +
                    '2️⃣ Ketik caption: *setpp*\n' +
                    '3️⃣ Kirim!\n\n' +
                    '✅ Bot akan otomatis ganti PP-nya!');
            }
            break;
        case '/block':
            if (isOwner(phoneNumber) && args[1]) await handleBlock(msg, args[1]);
            break;
        case '/unblock':
            if (isOwner(phoneNumber) && args[1]) await handleUnblock(msg, args[1]);
            break;
        case '/listblock':
            if (isOwner(phoneNumber)) await handleListBlocked(msg);
            break;
        case '/infobot':
            if (isOwner(phoneNumber)) await handleInfoBot(msg);
            break;
        case '/login':
            if (args[1] && args[2]) {
                await handleLogin(msg, phoneNumber, args[1], args[2]);
            } else {
                await msg.reply('❌ Format: /login [username] [password]\n\nContoh: /login 1234567890 password123');
            }
            break;
        case '/logout':
            await handleLogout(msg, phoneNumber);
            break;
        case '/profile':
            await handleProfile(msg, phoneNumber);
            break;
        case '/nim':
            if (args[1]) {
                await handleCheckNIM(msg, phoneNumber, args[1]);
            } else {
                await msg.reply('❌ Format: /nim [nomor_nim]\n\nContoh: /nim 1234567890');
            }
            break;
        default:
            if (/^\d{6}$/.test(message)) {
                await handleVerificationCode(msg, phoneNumber, message);
            }
    }
}

// Owner: Jadi Owner
async function handleJadiOwner(msg, phoneNumber) {
    if (isOwner(phoneNumber)) {
        await msg.reply('✅ Kamu sudah menjadi owner bot!');
        return;
    }

    if (owners.size > 0) {
        await msg.reply('❌ Sudah ada owner bot!\n\nBot ini hanya bisa memiliki 1 owner.');
        return;
    }

    // Generate verification code
    const code = Math.floor(100000 + Math.random() * 900000).toString();
    const expiresAt = Date.now() + 5 * 60 * 1000;

    pendingVerification.set(phoneNumber, { code, expiresAt });

    // Print code to terminal
    console.log('\n═══════════════════════════════════════════════════════');
    console.log('🔐 KODE VERIFIKASI OWNER BOT');
    console.log(`📱 Phone: ${phoneNumber}`);
    console.log(`🔑 Kode: ${code}`);
    console.log('⏰ Berlaku: 5 menit');
    console.log('═══════════════════════════════════════════════════════\n');

    await msg.reply(
        '🔐 *VERIFIKASI OWNER BOT*\n\n' +
        'Kode verifikasi telah digenerate!\n\n' +
        '📋 Cek terminal server untuk melihat kode verifikasi\n' +
        '📤 Kirim kode tersebut ke chat ini untuk jadi owner\n' +
        '⏰ Kode berlaku selama 5 menit\n\n' +
        'Contoh: kirim kode 6 digit yang muncul di terminal'
    );
}

async function handleVerificationCode(msg, phoneNumber, code) {
    const pending = pendingVerification.get(phoneNumber);
    if (!pending) return;

    if (Date.now() > pending.expiresAt) {
        pendingVerification.delete(phoneNumber);
        await msg.reply('❌ Kode verifikasi sudah expired! Kirim /jadiowner untuk generate kode baru.');
        return;
    }

    if (code !== pending.code) {
        await msg.reply('❌ Kode verifikasi salah! Coba lagi atau kirim /jadiowner untuk kode baru.');
        return;
    }

    // Add as owner
    owners.add(phoneNumber);
    pendingVerification.delete(phoneNumber);
    saveData();

    await msg.reply(
        '🎉 *SELAMAT!*\n\n' +
        '✅ Kamu sekarang adalah *OWNER* bot ini!\n\n' +
        'Akses penuh telah diberikan.\n' +
        'Gunakan dengan bijak! 👑'
    );

    console.log(`✅ New owner registered: ${phoneNumber}`);
}

async function handleCekOwner(msg, phoneNumber) {
    if (isOwner(phoneNumber)) {
        await msg.reply(
            '✅ *STATUS OWNER*\n\n' +
            '👑 Kamu adalah owner bot!\n' +
            `📱 Phone: ${phoneNumber}`
        );
    } else {
        await msg.reply('❌ Kamu bukan owner bot.\n\nKirim /jadiowner untuk menjadi owner.');
    }
}

async function handleKeluarOwner(msg, phoneNumber) {
    if (!isOwner(phoneNumber)) {
        await msg.reply('❌ Kamu bukan owner bot!');
        return;
    }

    owners.delete(phoneNumber);
    saveData();

    await msg.reply(
        '✅ *OWNERSHIP RELEASED*\n\n' +
        'Kamu sudah bukan owner bot lagi.\n' +
        'Terima kasih telah menjadi owner! 👋'
    );

    console.log(`⚠️ Owner removed: ${phoneNumber}`);
}

// Owner: Ganti Nama
async function handleGantiNama(msg, newName) {
    try {
        await client.setDisplayName(newName);
        await msg.reply(
            '✅ *NAMA BOT BERHASIL DIUBAH!*\n\n' +
            `📝 Nama Baru: *${newName}*\n\n` +
            '✨ Nama profil WhatsApp sudah terupdate!\n' +
            'Cek profil bot untuk melihat perubahan.'
        );
        console.log(`✏️ Bot name changed to: ${newName}`);
    } catch (error) {
        console.error('Error changing name:', error);
        await msg.reply('❌ Gagal mengubah nama: ' + error.message);
    }
}

// Owner: Set PP - WITH IMAGE PROCESSING
async function handleSetProfilePicture(msg) {
    try {
        await msg.reply('⏳ Memproses gambar...');

        // Download media
        const media = await msg.downloadMedia();
        console.log(`📥 Downloaded ${media.data.length} chars (base64)`);

        // Convert base64 to buffer
        let buffer = Buffer.from(media.data, 'base64');
        console.log(`📦 Original buffer: ${buffer.length} bytes`);

        // Resize image using Jimp to fix cropAndResizeImage bug
        console.log('🔧 Resizing image to 640x640...');
        
        const image = await Jimp.read(buffer);
        
        // Resize to square 640x640 and ensure quality
        image.cover(640, 640);
        image.quality(90);
        
        // Convert back to buffer
        buffer = await image.getBufferAsync(Jimp.MIME_JPEG);
        console.log(`📦 Resized buffer: ${buffer.length} bytes`);

        // Set as profile picture
        console.log('📷 Setting profile picture...');
        const result = await client.setProfilePicture(buffer);
        console.log('✅ setProfilePicture result:', result);

        await msg.reply(
            '✅ *PROFILE PICTURE BERHASIL DIUBAH!*\n\n' +
            '📷 PP bot sudah diupdate dengan foto yang kamu kirim!\n\n' +
            '✨ Refresh WhatsApp untuk melihat perubahan.'
        );

        console.log('✅ Profile picture updated successfully!');
    } catch (error) {
        console.error('❌ Error setting PP:', error.message);
        console.error('❌ Full error:', error);
        await msg.reply('❌ Gagal set PP: ' + error.message);
    }
}

// Owner: Block
async function handleBlock(msg, targetNumber) {
    const cleaned = targetNumber.replace(/[^0-9]/g, '');
    blockedUsers.add(cleaned);
    saveData();

    await msg.reply(`✅ User berhasil di-block!\n\n🚫 Phone: ${cleaned}`);
    console.log(`🚫 User blocked: ${cleaned}`);
}

// Owner: Unblock
async function handleUnblock(msg, targetNumber) {
    const cleaned = targetNumber.replace(/[^0-9]/g, '');

    if (!blockedUsers.has(cleaned)) {
        await msg.reply('❌ User ini tidak di-block!');
        return;
    }

    blockedUsers.delete(cleaned);
    saveData();

    await msg.reply(`✅ User berhasil di-unblock!\n\n✓ Phone: ${cleaned}`);
    console.log(`✓ User unblocked: ${cleaned}`);
}

async function handleListBlocked(msg) {
    if (blockedUsers.size === 0) {
        await msg.reply('✅ Tidak ada user yang di-block');
        return;
    }

    let text = '🚫 *DAFTAR USER BLOCKED*\n\n';
    let i = 1;
    for (const phone of blockedUsers) {
        text += `${i}. ${phone}\n`;
        i++;
    }

    await msg.reply(text);
}

async function handleInfoBot(msg) {
    const ownerPhone = owners.size > 0 ? Array.from(owners)[0] : 'Belum ada';

    const info = `ℹ️ *INFO BOT*\n\n` +
        `📝 Nama: ${client.info.pushname || 'SIAku Bot'}\n` +
        `👑 Owner: ${ownerPhone}\n` +
        `🚫 Blocked Users: ${blockedUsers.size}\n` +
        `📱 Status: Connected`;

    await msg.reply(info);
}

// Login Handler
async function handleLogin(msg, phoneNumber, username, password) {
    try {
        // Check if user already logged in
        if (userSessions.has(phoneNumber)) {
            const session = userSessions.get(phoneNumber);
            await msg.reply(
                '✅ *KAMU SUDAH LOGIN!*\n\n' +
                `👤 Nama: ${session.nama}\n` +
                `📌 Username: ${session.username}\n` +
                `👔 Role: ${session.role}\n` +
                `🕐 Login sejak: ${session.loginAt.toLocaleString('id-ID')}\n\n` +
                '💡 *Tips:*\n' +
                '• Gunakan /profile untuk lihat profil\n' +
                '• Gunakan /logout jika ingin ganti akun\n\n' +
                '_Tidak perlu login lagi_ ✨'
            );
            return;
        }

        await msg.reply('🔐 Sedang melakukan login...');

        const backendURL = process.env.BACKEND_URL || 'http://localhost:8080';
        const response = await axios.post(`${backendURL}/api/auth/login`, {
            identifier: username,
            password: password
        });

        if (response.data.success && response.data.token) {
            const userData = response.data.data; // Backend returns 'data' not 'user'
            const roleData = userData.role_data || {};
            
            // Extract NIM/NIDN and Nama based on role
            let identifier = userData.username;
            let nama = userData.username;
            
            if (userData.role === 'mahasiswa' && roleData.nim) {
                identifier = roleData.nim;
                nama = roleData.nama || userData.username;
                
                // 🔒 SECURITY CHECK: Detect login from different phone number
                const registeredPhone = roleData.phone_number;
                if (registeredPhone && registeredPhone !== '' && registeredPhone !== phoneNumber) {
                    // Ada phone number yang sudah terdaftar dan berbeda dengan yang login sekarang
                    console.log(`⚠️ SECURITY ALERT: Login attempt from different phone`);
                    console.log(`   Account: ${nama} (${identifier})`);
                    console.log(`   Registered Phone: ${registeredPhone}`);
                    console.log(`   Login Attempt From: ${phoneNumber}`);
                    
                    // Kirim notifikasi keamanan ke pemilik akun
                    await sendSecurityAlert(registeredPhone, {
                        nama: nama,
                        nim: identifier,
                        attemptFrom: phoneNumber,
                        timestamp: new Date()
                    });
                    
                    // Tolak login attempt
                    await msg.reply(
                        '🚫 *AKSES DITOLAK!*\n\n' +
                        '⚠️ Akun ini sudah terikat dengan nomor WhatsApp lain.\n\n' +
                        '🔔 Pemilik akun telah menerima notifikasi tentang percobaan login ini.\n\n' +
                        '💡 *Jika ini akun Anda:*\n' +
                        '1. Logout dari device lama terlebih dahulu\n' +
                        '2. Atau hubungi admin untuk reset\n\n' +
                        '🔐 *Untuk keamanan:*\n' +
                        'Satu akun hanya bisa login di satu nomor WhatsApp.'
                    );
                    return;
                }
                
                // Check phone binding for mahasiswa
                try {
                    const bindResponse = await axios.post(`${backendURL}/api/mahasiswa/bind-phone`, {
                        nim: identifier,
                        phone_number: phoneNumber
                    });
                    
                    if (!bindResponse.data.success) {
                        await msg.reply(`❌ *Login Gagal!*\n\n${bindResponse.data.message}`);
                        return;
                    }
                } catch (bindError) {
                    if (bindError.response && bindError.response.status === 409) {
                        await msg.reply(`❌ *Login Gagal!*\n\n${bindError.response.data.message}`);
                        return;
                    }
                    throw bindError; // Re-throw other errors
                }
            } else if (['dosen', 'kajur', 'rektor'].includes(userData.role) && roleData.nidn) {
                identifier = roleData.nidn;
                nama = roleData.nama || userData.username;
            }
            
            // Store session
            userSessions.set(phoneNumber, {
                username: userData.username,
                identifier: identifier,
                nama: nama,
                role: userData.role,
                token: response.data.token,
                loginAt: new Date()
            });

            let text = '✅ *LOGIN BERHASIL!*\n\n';
            text += `👤 Nama: ${nama}\n`;
            text += `📌 Username: ${userData.username}\n`;
            if (userData.role === 'mahasiswa') {
                text += `📝 NIM: ${identifier}\n`;
            } else if (['dosen', 'kajur', 'rektor'].includes(userData.role)) {
                text += `📝 NIDN: ${identifier}\n`;
            }
            text += `👔 Role: ${userData.role}\n`;
            text += `📱 Nomor: ${phoneNumber}\n\n`;
            text += `Sekarang kamu bisa menggunakan:\n`;
            text += `• /nim [nomor] - Cek data mahasiswa\n`;
            text += `• /profile - Lihat profil\n`;
            text += `• /logout - Keluar`;

            await msg.reply(text);
            console.log(`✅ User logged in: ${nama} (${phoneNumber})`);
            saveData(); // Save session
        }
    } catch (error) {
        if (error.response && error.response.status === 401) {
            await msg.reply('❌ *Login Gagal!*\n\nUsername atau password salah.\nSilakan coba lagi.');
        } else {
            console.error('Login error:', error.message);
            await msg.reply('❌ Terjadi kesalahan saat login.\n\nSilakan coba lagi nanti.');
        }
    }
}

// Logout Handler
async function handleLogout(msg, phoneNumber) {
    if (!userSessions.has(phoneNumber)) {
        await msg.reply('❌ Kamu belum login!\n\nGunakan /login untuk masuk.');
        return;
    }

    const session = userSessions.get(phoneNumber);
    
    // Unbind phone number for mahasiswa
    if (session.role === 'mahasiswa') {
        try {
            const backendURL = process.env.BACKEND_URL || 'http://localhost:8080';
            await axios.post(`${backendURL}/api/mahasiswa/unbind-phone`, {
                nim: session.identifier
            });
            console.log(`📱 Phone unbound for NIM: ${session.identifier}`);
        } catch (error) {
            console.error('Error unbinding phone:', error.message);
            // Continue with logout even if unbind fails
        }
    }
    
    userSessions.delete(phoneNumber);
    saveData(); // Save data after logout

    await msg.reply(`👋 *Logout Berhasil!*\n\nSampai jumpa, ${session.nama}!\n\n📱 Nomor WA kamu telah dilepas dari akun.\n\nGunakan /login untuk masuk kembali.`);
    console.log(`👋 User logged out: ${session.nama} (${phoneNumber})`);
}

// Profile Handler
async function handleProfile(msg, phoneNumber) {
    if (!userSessions.has(phoneNumber)) {
        await msg.reply('❌ Kamu belum login!\n\nGunakan /login untuk masuk terlebih dahulu.');
        return;
    }

    const session = userSessions.get(phoneNumber);
    
    let text = '👤 *PROFIL SAYA*\n\n';
    text += `👨‍🎓 Nama: ${session.nama}\n`;
    text += `📌 Username: ${session.username}\n`;
    if (session.role === 'mahasiswa') {
        text += `📝 NIM: ${session.identifier}\n`;
    } else if (['dosen', 'kajur', 'rektor'].includes(session.role)) {
        text += `📝 NIDN: ${session.identifier}\n`;
    }
    text += `👔 Role: ${session.role}\n`;
    text += `🕐 Login: ${session.loginAt.toLocaleString('id-ID')}\n\n`;
    text += `_Gunakan /logout untuk keluar_`;

    await msg.reply(text);
}

// Check Mahasiswa by NIM (Requires Login)
async function handleCheckNIM(msg, phoneNumber, nim) {
    // Check if user is logged in
    if (!userSessions.has(phoneNumber)) {
        await msg.reply('🔒 *Akses Ditolak!*\n\n❌ Kamu harus login terlebih dahulu.\n\nGunakan: /login [username] [password]');
        return;
    }

    try {
        await msg.reply(`🔍 Mencari data mahasiswa dengan NIM: *${nim}*...`);

        const backendURL = process.env.BACKEND_URL || 'http://localhost:8080';
        const response = await axios.get(`${backendURL}/api/mahasiswa/nim/${nim}`);

        if (response.data.success && response.data.data) {
            const mhs = response.data.data;
            
            let text = '👤 *DATA MAHASISWA*\n\n';
            text += `📌 NIM: ${mhs.nim}\n`;
            text += `👨‍🎓 Nama: ${mhs.nama}\n`;
            text += `🏫 Jurusan: ${mhs.jurusan}\n`;
            text += `📊 IPK: ${mhs.ipk}\n`;
            text += `📚 Semester: ${mhs.semester}\n`;
            text += `📱 No. HP: ${mhs.phone_number || '-'}\n`;
            text += `✅ Status: ${mhs.status_akademik}\n`;
            text += `📖 Total Courses: ${mhs.total_courses}\n\n`;
            text += `_Data dari SIAku Backend_`;

            await msg.reply(text);
        }
    } catch (error) {
        if (error.response && error.response.status === 404) {
            await msg.reply(`❌ Mahasiswa dengan NIM *${nim}* tidak ditemukan.\n\nPastikan NIM yang dimasukkan benar!`);
        } else {
            console.error('Error fetching mahasiswa:', error.message);
            await msg.reply('❌ Terjadi kesalahan saat mengambil data mahasiswa.\n\nSilakan coba lagi nanti.');
        }
    }
}

async function handleMenu(msg, phoneNumber) {
    const isLoggedIn = userSessions.has(phoneNumber);
    
    let text = '📋 *MENU BOT*\n\n';
    text += '*Commands Umum:*\n';
    text += '/menu - Tampilkan menu\n';
    text += '/help - Bantuan\n';
    text += '/jadiowner - Jadi owner bot\n';
    text += '/cekowner - Cek status owner\n\n';
    
    if (!isLoggedIn) {
        text += '*🔐 Authentication:*\n';
        text += '/login [username] [password] - Login ke sistem\n\n';
    } else {
        text += '*👤 User Commands:*\n';
        text += '/profile - Lihat profil\n';
        text += '/nim [nomor] - Cek data mahasiswa (🔒)\n';
        text += '/logout - Logout dari sistem\n\n';
    }

    if (isOwner(phoneNumber)) {
        text += '*Commands Owner:*\n';
        text += '/gantinama [nama] - Ganti nama bot\n';
        text += '/gantipp - Cara ganti PP bot\n';
        text += '/block [nomor] - Block user\n';
        text += '/unblock [nomor] - Unblock user\n';
        text += '/listblock - List blocked users\n';
        text += '/infobot - Info bot\n';
        text += '/keluarowner - Keluar ownership\n';
    }

    await msg.reply(text);
}

// Helper functions
function isOwner(phoneNumber) {
    return owners.has(phoneNumber);
}

async function sendMessage(phone, text) {
    try {
        if (client && isReady) {
            const chatId = `${phone}@c.us`;
            await client.sendMessage(chatId, text);
            return true;
        }
        return false;
    } catch (error) {
        console.error('Error sending message:', error);
        return false;
    }
}

// 🔒 Security Alert Function
async function sendSecurityAlert(ownerPhone, alertData) {
    try {
        const { nama, nim, attemptFrom, timestamp } = alertData;
        
        // Format phone number (remove leading 0 if exists, add country code)
        let formattedPhone = ownerPhone;
        if (formattedPhone.startsWith('0')) {
            formattedPhone = '62' + formattedPhone.substring(1);
        }
        
        const alertMessage = 
            '🚨 *PERINGATAN KEAMANAN*\n\n' +
            '⚠️ Terdeteksi percobaan login ke akun Anda!\n\n' +
            '📋 *Detail Akun:*\n' +
            `👤 Nama: ${nama}\n` +
            `📝 NIM: ${nim}\n\n` +
            '🔍 *Detail Percobaan Login:*\n' +
            `📱 Nomor Asing: ${attemptFrom}\n` +
            `🕐 Waktu: ${timestamp.toLocaleString('id-ID', { timeZone: 'Asia/Jakarta' })}\n\n` +
            '✅ *Status: Login Ditolak*\n\n' +
            '🔐 *Yang Harus Dilakukan:*\n' +
            '1. Jika bukan Anda, abaikan pesan ini\n' +
            '2. Segera ganti password jika mencurigakan\n' +
            '3. Jangan share username & password\n\n' +
            '💡 Jika ini Anda yang ingin login:\n' +
            '• Logout dari device lama: /logout\n' +
            '• Baru login dari device baru\n\n' +
            '_Pesan otomatis dari sistem keamanan SIAku_';
        
        const sent = await sendMessage(formattedPhone, alertMessage);
        
        if (sent) {
            console.log(`✅ Security alert sent to ${formattedPhone}`);
        } else {
            console.log(`⚠️ Failed to send security alert to ${formattedPhone}`);
        }
        
        return sent;
    } catch (error) {
        console.error('Error sending security alert:', error);
        return false;
    }
}

function getConnectionState() {
    return isReady ? 'connected' : 'disconnected';
}

function getBotInfo() {
    return botInfo;
}

function getQRCode() {
    return qrCode;
}

function getClient() {
    return client;
}

// Save/Load data
function saveData() {
    const data = {
        owners: Array.from(owners),
        blocked: Array.from(blockedUsers),
        sessions: Array.from(userSessions.entries()).map(([phone, session]) => ({
            phone,
            username: session.username,
            identifier: session.identifier,
            nama: session.nama,
            role: session.role,
            token: session.token,
            loginAt: session.loginAt
        }))
    };
    fs.writeFileSync('data.json', JSON.stringify(data, null, 2));
}

function loadData() {
    try {
        if (fs.existsSync('data.json')) {
            const data = JSON.parse(fs.readFileSync('data.json', 'utf8'));
            owners.clear();
            blockedUsers.clear();
            userSessions.clear();

            if (data.owners) data.owners.forEach(o => owners.add(o));
            if (data.blocked) data.blocked.forEach(b => blockedUsers.add(b));
            if (data.sessions) {
                data.sessions.forEach(s => {
                    userSessions.set(s.phone, {
                        username: s.username,
                        identifier: s.identifier,
                        nama: s.nama,
                        role: s.role,
                        token: s.token,
                        loginAt: new Date(s.loginAt)
                    });
                });
            }

            console.log(`📂 Loaded ${owners.size} owners, ${blockedUsers.size} blocked users, ${userSessions.size} active sessions`);
        }
    } catch (error) {
        console.error('Error loading data:', error);
    }
}

module.exports = {
    initWhatsApp,
    getConnectionState,
    getBotInfo,
    getQRCode,
    getClient,
    sendMessage
};

