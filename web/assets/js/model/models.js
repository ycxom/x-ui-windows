class User {

    constructor() {
        this.username = "";
        this.password = "";
    }
}

class Msg {

    constructor(success, msg, obj) {
        this.success = false;
        this.msg = "";
        this.obj = null;

        if (success != null) {
            this.success = success;
        }
        if (msg != null) {
            this.msg = msg;
        }
        if (obj != null) {
            this.obj = obj;
        }
    }
}

class DBInbound {

    constructor(data) {
        this.id = 0;
        this.userId = 0;
        this.up = 0;
        this.down = 0;
        this.total = 0;
        this.remark = "";
        this.enable = true;
        this.expiryTime = 0;

        this.listen = "";
        this.port = 0;
        this.protocol = "";
        this.settings = "";
        this.streamSettings = "";
        this.tag = "";
        this.sniffing = "";
        this.clientStats = "";
        
        // Backend address fields
        this.backendAddress = "";
        this.backendPort = 0;
        this.backendProtocol = "http";
        this.enableBackend = false;
        
        // 确保布尔值类型正确
        if (typeof this.enableBackend !== 'boolean') {
            this.enableBackend = Boolean(this.enableBackend);
        }
        if (data == null) {
            console.log('DBInbound: data is null');
            return;
        }
        
        console.log('🔥🔥🔥 DBInbound构造函数被调用 - 新版本 🔥🔥🔥');
        console.log('DBInbound: 原始数据 =', data);
        console.log('DBInbound: 后端代理字段检查:');
        console.log('  - enableBackend:', data.enableBackend);
        console.log('  - enable_backend:', data.enable_backend);
        console.log('  - backendAddress:', data.backendAddress);
        console.log('  - backend_address:', data.backend_address);
        console.log('  - backendPort:', data.backendPort);
        console.log('  - backend_port:', data.backend_port);
        console.log('  - backendProtocol:', data.backendProtocol);
        console.log('  - backend_protocol:', data.backend_protocol);
        
        // 在cloneProps之前先处理字段映射，确保数据正确复制
        // 手动处理数据库字段名到JavaScript属性名的映射
        if (data.backendProtocol !== undefined) {
            this.backendProtocol = data.backendProtocol;
        } else if (data.backend_protocol !== undefined) {
            this.backendProtocol = data.backend_protocol;
        }
        
        if (data.backendAddress !== undefined) {
            this.backendAddress = data.backendAddress;
        } else if (data.backend_address !== undefined) {
            this.backendAddress = data.backend_address;
        }
        
        if (data.backendPort !== undefined) {
            this.backendPort = Number(data.backendPort) || 0;
        } else if (data.backend_port !== undefined) {
            this.backendPort = Number(data.backend_port) || 0;
        }
        
        if (data.enableBackend !== undefined) {
            this.enableBackend = Boolean(data.enableBackend);
        } else if (data.enable_backend !== undefined) {
            // SQLite存储布尔值为数值，需要正确转换
            this.enableBackend = Boolean(Number(data.enable_backend));
        }
        
        console.log('DBInbound: 字段映射处理后(cloneProps前):');
        console.log('  - this.enableBackend:', this.enableBackend, typeof this.enableBackend);
        console.log('  - this.backendAddress:', this.backendAddress);
        console.log('  - this.backendPort:', this.backendPort);
        console.log('  - this.backendProtocol:', this.backendProtocol);
        
        // 现在调用cloneProps复制其他字段
        ObjectUtil.cloneProps(this, data);
        
        // 确保所有后端代理字段都有有效值
        if (this.backendProtocol === null || this.backendProtocol === undefined || this.backendProtocol === "") {
            this.backendProtocol = "http";
        }
        
        // 确保后端代理字段的响应式
        if (typeof Vue !== 'undefined' && Vue.set) {
            Vue.set(this, 'enableBackend', this.enableBackend);
            Vue.set(this, 'backendAddress', this.backendAddress);
            Vue.set(this, 'backendPort', this.backendPort);
            Vue.set(this, 'backendProtocol', this.backendProtocol);
        }
        
        console.log('DBInbound: 最终结果:');
        console.log('  - this.enableBackend:', this.enableBackend, typeof this.enableBackend);
        console.log('  - this.backendAddress:', this.backendAddress);
        console.log('  - this.backendPort:', this.backendPort);
        console.log('  - this.backendProtocol:', this.backendProtocol);
    }

    get totalGB() {
        return toFixed(this.total / ONE_GB, 2);
    }

    set totalGB(gb) {
        this.total = toFixed(gb * ONE_GB, 0);
    }

    get isVMess() {
        return this.protocol === Protocols.VMESS;
    }

    get isVLess() {
        return this.protocol === Protocols.VLESS;
    }

    get isTrojan() {
        return this.protocol === Protocols.TROJAN;
    }

    get isSS() {
        return this.protocol === Protocols.SHADOWSOCKS;
    }

    get isSocks() {
        return this.protocol === Protocols.SOCKS;
    }

    get isHTTP() {
        return this.protocol === Protocols.HTTP;
    }

    get address() {
        let address = location.hostname;
        if (!ObjectUtil.isEmpty(this.listen) && this.listen !== "0.0.0.0") {
            address = this.listen;
        }
        return address;
    }

    get _expiryTime() {
        if (this.expiryTime === 0) {
            return null;
        }
        return moment(this.expiryTime);
    }

    set _expiryTime(t) {
        if (t == null) {
            this.expiryTime = 0;
        } else {
            this.expiryTime = t.valueOf();
        }
    }

    get isExpiry() {
        return this.expiryTime < new Date().getTime();
    }

    toInbound() {
        let settings = {};
        if (!ObjectUtil.isEmpty(this.settings)) {
            settings = JSON.parse(this.settings);
        }

        let streamSettings = {};
        if (!ObjectUtil.isEmpty(this.streamSettings)) {
            streamSettings = JSON.parse(this.streamSettings);
        }

        let sniffing = {};
        if (!ObjectUtil.isEmpty(this.sniffing)) {
            sniffing = JSON.parse(this.sniffing);
        }

        const config = {
            port: this.port,
            listen: this.listen,
            protocol: this.protocol,
            settings: settings,
            streamSettings: streamSettings,
            tag: this.tag,
            sniffing: sniffing,
            clientStats: this.clientStats,
        };
        return Inbound.fromJson(config);
    }

    hasLink() {
        switch (this.protocol) {
            case Protocols.VMESS:
            case Protocols.VLESS:
            case Protocols.TROJAN:
            case Protocols.SHADOWSOCKS:
                return true;
            default:
                return false;
        }
    }

    genLink() {
        const inbound = this.toInbound();
        return inbound.genLink(this.address, this.remark);
    }
}

class AllSetting {

    constructor(data) {
        this.webListen = "";
        this.webPort = 54321;
        this.webCertFile = "";
        this.webKeyFile = "";
        this.webBasePath = "/";
        this.tgBotEnable = false;
        this.tgBotToken = "";
        this.tgBotChatId = 0;
        this.tgRunTime = "";
        this.xrayTemplateConfig = "";

        this.timeLocation = "Asia/Shanghai";

        if (data == null) {
            return
        }
        ObjectUtil.cloneProps(this, data);
    }

    equals(other) {
        return ObjectUtil.equals(this, other);
    }
}