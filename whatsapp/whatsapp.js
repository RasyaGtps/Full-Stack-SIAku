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
        case '/nim':
            if (args[1]) {
                await handleCheckNIM(msg, args[1]);
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

// Check Mahasiswa by NIM
async function handleCheckNIM(msg, nim) {
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
    let text = '📋 *MENU BOT*\n\n';
    text += '*Commands Umum:*\n';
    text += '/menu - Tampilkan menu\n';
    text += '/help - Bantuan\n';
    text += '/nim [nomor] - Cek data mahasiswa\n';
    text += '/jadiowner - Jadi owner bot\n';
    text += '/cekowner - Cek status owner\n\n';

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
        blocked: Array.from(blockedUsers)
    };
    fs.writeFileSync('data.json', JSON.stringify(data, null, 2));
}

function loadData() {
    try {
        if (fs.existsSync('data.json')) {
            const data = JSON.parse(fs.readFileSync('data.json', 'utf8'));
            owners.clear();
            blockedUsers.clear();

            if (data.owners) data.owners.forEach(o => owners.add(o));
            if (data.blocked) data.blocked.forEach(b => blockedUsers.add(b));

            console.log(`📂 Loaded ${owners.size} owners and ${blockedUsers.size} blocked users`);
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

