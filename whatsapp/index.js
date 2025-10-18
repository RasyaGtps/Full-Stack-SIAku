const express = require('express');
const bodyParser = require('body-parser');
const { initWhatsApp, getConnectionState, getBotInfo } = require('./whatsapp');
const routes = require('./routes');

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

// CORS
app.use((req, res, next) => {
    res.header('Access-Control-Allow-Origin', '*');
    res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept, Authorization');
    res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
    next();
});

// Routes
app.use('/api/wa', routes);

// Health check
app.get('/', (req, res) => {
    const isConnected = getConnectionState() === 'connected';
    const info = getBotInfo();
    res.json({ 
        success: true,
        service: 'WhatsApp Bot Service',
        status: isConnected ? 'connected' : 'disconnected',
        connected: isConnected,
        bot_name: info.name,
        bot_number: info.number,
        port: PORT
    });
});

// Start server
app.listen(PORT, async () => {
    console.log(`\n🚀 WhatsApp Service running on port ${PORT}`);
    console.log(`📱 Initializing WhatsApp connection...\n`);
    
    // Initialize WhatsApp
    await initWhatsApp();
});

// Graceful shutdown
process.on('SIGINT', () => {
    console.log('\n⚠️  Shutting down gracefully...');
    process.exit(0);
});

