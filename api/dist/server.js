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
var compression_1 = __importDefault(require("compression"));
var cors_1 = __importDefault(require("cors"));
var dotenv = __importStar(require("dotenv"));
var express_1 = __importDefault(require("express"));
var helmet_1 = __importDefault(require("helmet"));
var http_1 = require("http");
dotenv.config();
console.log(dotenv.load());
console.log('moi');
console.log(process.env.API_PORT);
var port = process.env.API_PORT;
var app = express_1.default();
app.use(helmet_1.default());
app.disable('x-powered-by');
app.use('*', cors_1.default());
app.use(compression_1.default());
app.get('/', function (_req, res) { return res.send('Hello World!'); });
var httpServer = http_1.createServer(app);
httpServer.listen({ port: port }, function () {
    console.log("\uD83D\uDE80 Test api is running on port " + port);
});
