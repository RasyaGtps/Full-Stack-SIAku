const express = require('express');
const router = express.Router();
const { getConnectionState, getQRCode, sendMessage, getSock } = require('./whatsapp');

// Get connection status
router.get('/status', (req, res) => {
    res.json({
        success: true,
        data: {
            connected: getConnectionState() === 'open',
            status: getConnectionState()
        }
    });
});

// Get QR Code
router.get('/qr-code', (req, res) => {
    const qr = getQRCode();
    
    if (!qr) {
        return res.json({
            success: false,
            message: 'No QR code available. Already connected or reconnecting.'
        });
    }

    res.json({
        success: true,
        data: {
            qr_code: qr,
            message: 'Scan this QR code with WhatsApp'
        }
    });
});

// Send message
router.post('/send', async (req, res) => {
    const { phone_number, message } = req.body;

    if (!phone_number || !message) {
        return res.status(400).json({
            success: false,
            error: 'phone_number and message are required'
        });
    }

    try {
        await sendMessage(phone_number, message);
        res.json({
            success: true,
            message: 'Message sent successfully'
        });
    } catch (error) {
        res.status(500).json({
            success: false,
            error: error.message
        });
    }
});

// Broadcast message
router.post('/broadcast', async (req, res) => {
    const { phone_numbers, message } = req.body;

    if (!phone_numbers || !Array.isArray(phone_numbers) || !message) {
        return res.status(400).json({
            success: false,
            error: 'phone_numbers (array) and message are required'
        });
    }

    try {
        for (const phone of phone_numbers) {
            await sendMessage(phone, message);
        }

        res.json({
            success: true,
            message: `Broadcast sent to ${phone_numbers.length} recipients`
        });
    } catch (error) {
        res.status(500).json({
            success: false,
            error: error.message
        });
    }
});

module.exports = router;

