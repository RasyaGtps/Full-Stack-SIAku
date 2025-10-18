const fs = require('fs');

console.log('\n🗑️  Clearing WhatsApp session...\n');

// Remove auth_session folder (whatsapp-web.js)
if (fs.existsSync('auth_session')) {
    fs.rmSync('auth_session', { recursive: true, force: true });
    console.log('✅ Removed: auth_session/');
}

// Remove auth_info folder (baileys)
if (fs.existsSync('auth_info')) {
    fs.rmSync('auth_info', { recursive: true, force: true });
    console.log('✅ Removed: auth_info/');
}

// Remove .wwebjs_auth folder
if (fs.existsSync('.wwebjs_auth')) {
    fs.rmSync('.wwebjs_auth', { recursive: true, force: true });
    console.log('✅ Removed: .wwebjs_auth/');
}

// Remove data file
if (fs.existsSync('data.json')) {
    fs.unlinkSync('data.json');
    console.log('✅ Removed: data.json');
}

console.log('\n✅ Session cleared! Run "npm start" to scan new QR code.\n');

