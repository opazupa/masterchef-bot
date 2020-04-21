"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) result[k] = mod[k];
    result["default"] = mod;
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
var express_1 = __importDefault(require("express"));
var http_1 = require("http");
var compression_1 = __importDefault(require("compression"));
var cors_1 = __importDefault(require("cors"));
var dotenv = __importStar(require("dotenv"));
var helmet_1 = __importDefault(require("helmet"));
var port = process.env.PORT || 3000;
var app = express_1.default();
dotenv.config();
app.use(helmet_1.default());
app.use("*", cors_1.default());
app.use(compression_1.default());
app.get("/", function (_req, res) { return res.send("Hello World!"); });
var httpServer = http_1.createServer(app);
httpServer.listen({ port: port }, function () {
    return console.log("\uD83D\uDE80 Test api is running on port " + port);
});
