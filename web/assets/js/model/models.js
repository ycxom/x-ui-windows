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
        if (data == null) {
            return;
        }
        ObjectUtil.cloneProps(this, data);
        
        // 手动处理数据库字段名到JavaScript属性名的映射
        if (data.backend_protocol !== undefined) {
            this.backendProtocol = data.backend_protocol;
        }
        if (data.backend_address !== undefined) {
            this.backendAddress = data.backend_address;
        }
        if (data.backend_port !== undefined) {
            this.backendPort = data.backend_port;
        }
        if (data.enable_backend !== undefined) {
            this.enableBackend = Boolean(data.enable_backend);
        }
        
        // 确保所有后端代理字段都有有效值
        if (this.backendProtocol === null || this.backendProtocol === undefined || this.backendProtocol === "") {
            this.backendProtocol = "http";
        }
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