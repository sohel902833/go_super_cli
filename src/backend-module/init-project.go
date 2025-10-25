package backendmodule

import (
	"github.com/sohel902833/go_super_cli/src/types"
)


func GetInitProjectInstructions()([]types.FileInstruction,[]types.UpdateInstruction) {
	var fileInstructions= []types.FileInstruction{
		{
			 FilePath: "./tsconfig.json",
			 Description: "Creating tsconfig.json",
			 Content:`{
				"compileOnSave": false,
				"compilerOptions": {
					"target": "ESNext",
					"lib": ["ES6"],
					"allowJs": true,
					"module": "CommonJS",
					"rootDir": ".",
					"outDir": "./dist",
					"esModuleInterop": true,
					"strict": true,
					"skipLibCheck": true,
					"forceConsistentCasingInFileNames": true,
					"moduleResolution": "node",
					"resolveJsonModule": true,
					"allowSyntheticDefaultImports": true,
					"typeRoots": ["./src/types", "./node_modules/@types"],
					"sourceMap": true,
					"types": ["node", "express"],
					"noImplicitAny": false,
					"baseUrl": "./src",
					"paths": {
						"@/*": ["*"]
					}
				},
				"include": ["src/**/*"],
				"exclude": ["node_modules"]
			}`,
		},
		{
			FilePath: "./package.json",
			Description: "Creating package json file",
			Content: `{
    "name": "{{PROJECT_NAME}}",
    "version": "1.0.0",
    "description": "",
    "main": "index.js",
    "scripts": {
        "test": "echo \"Error: no test specified\" && exit 1",
        "dev": "ts-node-dev -r tsconfig-paths/register ./src/index.ts",
        "build": "tsc && tsc-alias",
        "build-permission": "ts-node-dev -r tsconfig-paths/register ./src/app/role/permission-creator.ts"
    },
    "keywords": [],
    "author": "",
    "license": "ISC",
    "dependencies": {
        "bcrypt": "^5.1.1",
        "cloudinary": "^2.6.0",
        "cookie-parser": "^1.4.7",
        "cors": "^2.8.5",
        "dotenv": "^16.4.5",
        "express": "^4.21.1",
        "jsonwebtoken": "^9.0.2",
        "mongoose": "^7.5.2",
        "multer": "^1.4.5-lts.2",
        "nodemailer": "^6.10.0",
        "sharp": "^0.34.1",
        "validator": "^13.12.0",
        "zod": "^3.23.8"
    },
    "devDependencies": {
        "@types/cookie-parser": "^1.4.7",
        "@types/cors": "^2.8.17",
        "@types/express": "^5.0.0",
        "@types/node": "^22.8.1",
        "ts-node-dev": "^2.0.0",
        "tsc": "^2.0.4",
        "tsc-alias": "^1.8.10",
        "tsconfig-paths": "^4.2.0",
        "typescript": "^5.6.3"
    }
}
`,
		},
		{
			FilePath: "src/index.ts",
			Description: "Creating Project Root File",
			Content: `import { Server, createServer } from "http";
import app from "@/app";
import config from "@/config";
import * as db from "@/db";
import { IUsers } from "./app/users/users.types";
import { IUserPermission } from "./app/role/role.types";
import { Cloudinary } from "./helpers/cloudinary";
let server: Server;

declare global {
    namespace Express {
        export interface Request {
            userId?: string;
            user?: IUsers;
            userPermissions?: IUserPermission;
            files: any;
        }
    }
}
async function startServer() {
    try {
        db.connect();
        Cloudinary();
        server = createServer(app);
        server.listen(config.port, () => {
            console.log(__BACKTICK__Application is Running On Port: ${config.port}__BACKTICK__);
        });
    } catch (err) {
        console.log("Failed to connect database", err);
    }
}

startServer();
`,
		},
		{
			FilePath: "src/app.ts",
			Description: "Creating app file",
			Content: `import express, { Application } from "express";
import cors from "cors";
import cookieParser from "cookie-parser";
import appRoutes from "@/app/index";
import globalErrorHandler from "./middlewares/globalErrorHandler";
import path from "path";
const app: Application = express();

const allowedOrigins = ["http://localhost:5173"];

//setup cors here
// app.use(cors());
app.use(
    cors({
        origin: (origin, callback) => {
            if (!origin || allowedOrigins.includes(origin)) {
                callback(null, true);
            } else {
                callback(new Error("Not allowed by CORS"));
            }
        },
        credentials: true,
    })
);

// parser
app.use(cookieParser());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

app.use("/uploads/", express.static(path.join(__dirname, "../../", "uploads")));
// Application router

app.use("/api/v1/", appRoutes);

app.get("/health", (_req, res) => {
    res.status(200).json({
        status: "UP",
    });
});

//404 handler
app.use((_req, res) => {
    res.status(404).json({
        message: "Not Found",
    });
});

// error handling
app.use(globalErrorHandler);
app.use((err: any, _req: any, res: any, _next: any) => {
    console.log(err.stack);
    res.status(500).json({
        message: "Internal server error",
    });
});

export default app;
`,
		},


		{
			FilePath:    "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.controller.ts",
			Description: "Creating controller with CRUD operations",
			Content: `import { Request, Response } from 'express';
import { {{PASCAL_CASE_MODULE_NAME}}Service } from './{{LOWER_CASE_MODULE_NAME}}.service';

export class {{PASCAL_CASE_MODULE_NAME}}Controller {
  private service: {{PASCAL_CASE_MODULE_NAME}}Service;

  constructor() {
    this.service = new {{PASCAL_CASE_MODULE_NAME}}Service();
  }

  async create(req: Request, res: Response) {
    try {
      const data = await this.service.create(req.body);
      res.status(201).json(data);
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  }

  async findAll(req: Request, res: Response) {
    try {
      const data = await this.service.findAll();
      res.status(200).json(data);
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  }

  async findOne(req: Request, res: Response) {
    try {
      const data = await this.service.findOne(req.params.id);
      res.status(200).json(data);
    } catch (error) {
      res.status(404).json({ error: error.message });
    }
  }

  async update(req: Request, res: Response) {
    try {
      const data = await this.service.update(req.params.id, req.body);
      res.status(200).json(data);
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  }

  async delete(req: Request, res: Response) {
    try {
      await this.service.delete(req.params.id);
      res.status(204).send();
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  }
}
`,
		},
		{
			FilePath:    "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.service.ts",
			Description: "Creating service layer with business logic",
			Content: `import { {{PASCAL_CASE_MODULE_NAME}}Model } from './{{LOWER_CASE_MODULE_NAME}}.model';
import { {{PASCAL_CASE_MODULE_NAME}}Schema } from './{{LOWER_CASE_MODULE_NAME}}.schema';

export class {{PASCAL_CASE_MODULE_NAME}}Service {
  async create(data: any) {
    const validated = {{PASCAL_CASE_MODULE_NAME}}Schema.parse(data);
    return await {{PASCAL_CASE_MODULE_NAME}}Model.create(validated);
  }

  async findAll() {
    return await {{PASCAL_CASE_MODULE_NAME}}Model.find();
  }

  async findOne(id: string) {
    const record = await {{PASCAL_CASE_MODULE_NAME}}Model.findById(id);
    if (!record) {
      throw new Error('{{PASCAL_CASE_MODULE_NAME}} not found');
    }
    return record;
  }

  async update(id: string, data: any) {
    const validated = {{PASCAL_CASE_MODULE_NAME}}Schema.partial().parse(data);
    return await {{PASCAL_CASE_MODULE_NAME}}Model.findByIdAndUpdate(id, validated, { new: true });
  }

  async delete(id: string) {
    return await {{PASCAL_CASE_MODULE_NAME}}Model.findByIdAndDelete(id);
  }
}
`,
		},
		{
			FilePath: "src/config/index.ts",
			Description: "Creating config file",
			Content: `import dotenv from "dotenv";
import path from "path";
dotenv.config({ path: path.join(process.cwd(), ".env") });

export default {
    env: process.env.NODE_ENV as "development" | "stage" | "production",
    port: process.env.PORT,
    database_url: process.env.DATABASE_URL,
    default_user_pass: process.env.DEFAULT_USER_PASS,
    jwt_expiry: process.env.JWT_EXPIRY,
    jwt_secret: process.env.USER_JWT_SECRET,
    refresh_token_secret: process.env.USER_REFRESH_TOKEN_SECRET,
    refresh_token_expiry: Number(process.env.JWT_EXPIRY) * 4,
    verify_email_secret: process.env.VERIFY_EMAIL_SECRET,
    verify_email_expiry: process.env.VERIFY_EMAIL_EXPIRY_TIME,
    mailer_user_name: process.env.MAILER_USER_NAME,
    mailer_user_password: process.env.MAILER_USER_PASSWORD,
    frontend_url: process.env.FRONTEND_URL,
    cloudinary: {
        name: process.env.CLOUDINARY_NAME,
        api: process.env.CLOUDINARY_API_KEY,
        secret: process.env.CLOUDINARY_API_SECRET,
    },
};
`,

		},
		{
			FilePath: "src/db/index.ts",
			Description: "Creating db init file",
			Content: `import config from "@/config";
import mongoose from "mongoose";
export { default as MODEL_NAMES } from "./modelNames";
export { default as models } from "./models";
export const connect = async () => {
    try {
        await mongoose.connect(config.database_url as string);
        console.log("Database Successfully Connected");
    } catch (err) {
        console.log("Database Failed To Connect", err);
    }
};
`},
		{
			FilePath: "src/modelNames.ts",
			Description: "Creating model names files",
			Content: `const MODEL_NAMES = {
    //MODEL_NAME_DEFINATION_AREA
};

export default MODEL_NAMES;
`,
		},
		{
			FilePath: "src/models.ts",
			Description: "Creating models listing files",
			Content: `import Role from "@/app/role/role.model";
import Users, { VerifyCode } from "@/app/users/users.model";
//MODEL_IMPORT_DEFINATION_AREA

const models = {
    Users,
    Role,
    VerifyCode,
	//MODEL_NAME_DEFINE_AREA
};

export default models;`,
		},
		{
			FilePath:"src/errors/index.ts",
			Description: "Creating error creator file",
			Content: `export const createError = (
    error: string,
    opts?: {
        errorMessages: Record<string, any>;
    }
) => {
    const { errorMessages = {} } = opts || {};
    const newError = new Error(error);
    //@ts-ignore
    newError.errorMessages = errorMessages;
    throw newError;
};
`,
		},
		{
			FilePath: "src/helpers/cloudinary.ts",
			Description: "Cloudinary helpers",
			Content: `import config from "@/config";
import { v2 as cloudinary } from "cloudinary";

export const Cloudinary = async () => {
    await cloudinary.config({
        cloud_name: config?.cloudinary?.name,
        api_key: config?.cloudinary?.api,
        api_secret: config?.cloudinary?.secret,
    });
};
`,
		},
		{
			FilePath:"src/index.ts",
			Description: "Creating helpers file",
			Content: `import mongoose from "mongoose";

const isObjectId = (val: string) => {
    return mongoose.Types.ObjectId.isValid(val);
};
const modifyValue = (val: string) => {
    if (val === "true" || val === "false") {
        return val === "true";
    } else if (isObjectId(val)) {
        return val;
    }
    return { $regex: val, $options: "i" };
};

export const modifyQuery = (queryInput: { [key: string]: string }) => {
    const query = { ...queryInput };
    let page = 1;
    let limit = 10;
    const finalQuery: Record<string, any> = {};
    const sort: Record<string, 1 | -1> = {};

    if (query?.page) {
        page = Number(query?.page);
        delete query.page;
    }

    if (query?.limit) {
        limit = Number(query?.limit);
        delete query.limit;
    }
    Object.keys(query).forEach((key) => {
        if (key.startsWith("_sort_")) {
            const field = key.replace("_sort_", "");
            const value = query[key]?.toLowerCase();
            sort[field] = value === "desc" ? -1 : 1;
            delete query[key];
        }
    });
    if (!Object.keys(sort)?.length) {
        sort.createdAt = -1;
    }

    Object.keys(query).forEach((queryKey) => {
        const queryValue = query[queryKey];
        if (Array.isArray(queryValue)) {
            if (!finalQuery?.$or) {
                finalQuery.$or = [];
            }
            queryValue.forEach((item) => {
                finalQuery.$or.push({
                    [queryKey]: modifyValue(item),
                });
            });
        } else if (typeof queryValue === "string" && queryValue.includes(",")) {
            const splitedList = queryValue.split(",");
            if (!finalQuery?.$or) {
                finalQuery.$or = [];
            }
            splitedList.forEach((item) => {
                finalQuery.$or.push({
                    [queryKey]: modifyValue(item),
                });
            });
        } else {
            finalQuery[queryKey] = modifyValue(queryValue);
        }
    });

    return {
        page,
        limit,
        finalQuery,
        sort,
    };
};
`,
		},
		{
			FilePath: "src/helpers/jwtHelpers.ts",
			Description: "Creating jwt helpers file",
			Content: `import { sign, verify } from "jsonwebtoken";
import config from "../config";
export const getSignedToken = (userId: string): string => {
  const token = sign({ userId }, config.jwt_secret as string, {
    expiresIn: config.jwt_expiry,
  });
  return token;
};

export const getRefreshToken = (userId: string): string => {
  const token = sign({ userId }, config.refresh_token_secret as string, {
    expiresIn: config.refresh_token_expiry,
  });
  return token;
};

export const verifyRefreshToken = async (token: string) => {
  try {
    const decoded: any = verify(token, config.refresh_token_secret as string);
    if (!decoded) {
      return {
        success: false,
        message: "Refresh token expired",
        errorFor: "auth",
      };
    }
    return {
      success: true,
      userId: decoded.userId,
    };
  } catch (err) {
    return {
      success: false,
      message: "Refresh token expired",
      errorFor: "auth",
    };
  }
};

export const getVerifyMailToken = (userId: string): string => {
  const token = sign({ userId }, config.verify_email_secret as string, {
    expiresIn: config.verify_email_expiry,
  });
  return token;
};
`,
		},
		{
			FilePath: "src/helpers/mailer.ts",
			Description: "Creating mailers file",
			Content: `import {
    createTransport,
    createTestAccount,
    getTestMessageUrl,
} from "nodemailer";
import configEnv from "../config/index";
interface IEmail {
    to: string;
    subject: string;
    html?: string;
}

export const sendMail = async ({ to, html, subject }: IEmail) => {
    const userName = configEnv.mailer_user_name;
    const password = configEnv.mailer_user_password;
    const transporter = createTransport({
        host: "smtp.gmail.com",
        port: 465,
        secure: true,
        auth: {
            user: userName,
            pass: password,
        },
    });
    const info = await transporter.sendMail({
        from: userName,
        to: to,
        subject: subject,
        html: html,
    });

    // console.log("Mail send", info);
};
`,
		},
		{
			FilePath: "src/helpers/mailTemplate.ts",
			Description: "Creating mail template file",
			Content: `export const getVerifyMailTemplate = (url: string, text?: string): string => {
  return __BACKTICK__<a href='www.google.com'>Click Here to verify your email</a>__BACKTICK__;
}
export const getPasswordResetCodeTemplate = (
  code: number,
  text?: string
): string => {
  return __BACKTICK__<div
  style="
    text-align: center;
    border-radius: 5px;
    padding: 10px;
    border: 1px solid blue;
    height: 600px;
    width: 600px;
    background-color: white;
    margin: 0 auto;
  "
>
  <h1
    style="
      text-align: center;
      color: blue;
      font-weight: 600;
      font-family: Arial, Helvetica, sans-serif;
    "
  >
    Use Below Code for reset your password
  </h1>
  <button
    style="
      outline: none;
      border: none;
      font-weight: bolder;
      cursor: pointer;
      color: white;
      background-color: blue;
      padding: 10px 20px;
      margin-top: 40px;
      font-size: 30px;
    "
  >
    ${code}
  </button>
  <p style="margin-top: 20px">${text}</p>
</div>__BACKTICK__;
};`},
			{
				FilePath: "src/helpers/pagination.ts",
				Description: "Creating pagination helpers",
				Content: `import { Model, Query } from "mongoose";

export interface IPaginationResult {
    pagination: {
        total: number;
        next_page: number;
        prev_page: number;
        limit: number;
    };
    data: null | any[];
}
export interface IPaginationReturnVal {
    result: IPaginationResult;
    limit: number;
    startIndex: number;
}

export const getPaginationProperty = async (
    page: number,
    limit: number,
    model: Model<any>,
    filter: any
): Promise<IPaginationReturnVal> => {
    const startIndex = (page - 1) * limit;
    const endIndex = page * limit;
    let result: IPaginationResult = {
        pagination: {
            limit: 0,
            next_page: 0,
            prev_page: 0,
            total: 0,
        },
        data: [],
    };

    const totalDocuments = await model.countDocuments(filter).exec();
    result.pagination.total = totalDocuments;
    if (endIndex < totalDocuments) {
        result.pagination.next_page = page + 1;
        result.pagination.limit = limit;
    }
    if (startIndex > 0) {
        result.pagination.prev_page = page - 1;
        result.pagination.limit = limit;
    }

    return { result, limit: limit, startIndex };
};

export const getWithPagination = async ({
    page,
    limit,
    model,
    filter = {},
    populate,
    projection = null,
    sort = { createdAt: -1 },
}: {
    page: number;
    limit: number;
    model: Model<any>;
    filter?: any;
    populate?: any;
    projection?: any;
    sort?: any;
}) => {
    let {
        result,
        limit: lm,
        startIndex,
    } = await getPaginationProperty(page, limit, model, filter);

    let query = model
        .find(filter, projection)
        .skip(startIndex)
        .limit(lm)
        .sort(sort);
    if (populate) {
        query = query.populate(populate);
    }
    const data = await query.exec();

    result.data = data;
    return result;
};
`,
			},
			{	
				FilePath: "src/helpers/randomNumber.ts",
				Description: "Random number file",
				Content: `export const getRandomNumber = (min?: number, max?: number) => {
    const minm = min ? min : 100000;
    const maxm = max ? max : 999999;
    return Math.floor(Math.random() * (maxm - minm + 1)) + minm;
};
`,

			},
			{
				FilePath: "src/middlewares/authGard.ts",
				Description: "Creating auth gard",
				Content: `import { Permissions } from "@/app/role";
import { IUsers } from "@/app/users/users.types";
import config from "@/config";
import { NextFunction, Response, Request } from "express";
import { verify } from "jsonwebtoken";
import db from "@/db/models";
import { hasPermissions } from "@/permissions";
import { IRole, IUserPermission } from "@/app/role/role.types";
const { Users } = db;

export const authGard = (
    permissions?: Permissions[],
    opts?: {
        needToActivate?: boolean;
        addUser?: boolean;
    }
) => {
    return async (
        req: Request,
        res: Response,
        next: NextFunction
    ): Promise<any> => {
        try {
            const { needToActivate, addUser } = opts || {};
            let token = req.headers.authorization || req.cookies.Authorization;
            const refreshToken =
                req.headers.refreshToken || req.cookies.refreshToken;
            //check token exists or not
            if (!token) {
                return res.status(404).json({
                    message: "Authentication Error.",
                    errorFor: "auth",
                });
            }
            //check token verified or not
            const decoded: any = verify(token, config.jwt_secret as string);
            const dbUser: IUsers | null = await Users.findById({
                _id: decoded.userId,
            }).populate("role");
            //@ts-ignore
            const userRole: IRole = dbUser?.role;
            //@ts-ignore
            const userPermissions: IUserPermissionn = userRole?.permissions;
            if (!dbUser) {
                return res.status(200).json({
                    message: "Requested User Was Not Found.",
                });
            }

            if (needToActivate && !dbUser?.activated) {
                return res.status(200).json({
                    message: "Your Account Temporary Deactivated.",
                });
            }
            const hasPerm = hasPermissions(permissions, userPermissions);
            // console.log("has Perm", hasPerm, { permissions, userPermissions });
            if (permissions?.length && !hasPerm) {
                return res.status(401).json({
                    message: "You are Not Authorized to Perform this Action",
                    type: "unauthorized",
                });
            }
            // @ts-ignore
            const newUser = { ...dbUser._doc };
            delete newUser.password;
            req.userId = decoded.userId;
            req.user = newUser;
            req.userPermissions = userPermissions;

            next();
        } catch (err: any) {
            console.log("Error", err);
            return res.status(404).json({
                message: err.message,
                errorFor: "auth",
            });
        }
    };
};
`,
			},
			{
				FilePath: "src/middlewares/globalErrorHandler.ts",
				Description: "Creating global error handler file",
				Content:`import { ZodError } from "zod";
import { Request, Response, NextFunction, ErrorRequestHandler } from "express";
import { MongooseError } from "mongoose";
const generateErrorObject = (issues: any[] = []) => {
    let errors = {};
    if (issues && Array.isArray(issues)) {
        issues.forEach((issue) => {
            const message = issue?.message;
            let currentLevel = errors;
            const path: any[] = issue?.path ?? [];
            if (path?.length === 2 && path[0] === "body") {
                const pathName = path[path.length - 1];
                errors[pathName] = message;
            } else {
                path.forEach((key, index) => {
                    if (key === "body") {
                        return;
                    } else if (index === path.length - 1) {
                        currentLevel[key] = message;
                    } else {
                        currentLevel[key] = currentLevel[key] || {};
                        currentLevel = currentLevel[key];
                    }
                });
            }
        });
    }
    return errors;
};

const handleZodError = (error: any) => {
    // let errors = {}
    const errors = generateErrorObject(error?.issues ?? []);

    const statusCode = 400;

    return {
        statusCode,
        message: "Validation Error",
        errorMessages: errors,
    };
};

const globalErrorHandler = (
    err: any,
    req: Request,
    res: Response,
    next: NextFunction
): any => {
    try {
        let statusCode = 500;
        let message = "Something went wrong!";
        let errorMessages;
        const errorKey = err?.errorKey;
        if (err instanceof ZodError) {
            const simplefiedError = handleZodError(err);
            statusCode = simplefiedError.statusCode;
            message = simplefiedError.message;
            errorMessages = simplefiedError.errorMessages;
        } else if (err instanceof Error) {
            // console.log("Mongoose Error", err);
            message = err?.message;
            //@ts-ignore
            errorMessages = err?.errorMessages ? err?.errorMessages : [];
        } else {
            return res.status(404).json({
                message: "Server Error Found",
                error: err,
            });
        }
        const result: any = {
            success: false,
            message,
            errorMessages,
            // stack: err?.stack,
        };
        if (errorKey === "already_exists") {
            result.alreadyExists = true;
        }
        res.status(statusCode).send(result);
    } catch (err) {
        console.log("Err", err);
        return res.status(404).json({
            message: "Server Error Found",
            error: err,
        });
    }
};

export default globalErrorHandler;
`,
			},
			{
				FilePath: "src/app/index.ts",
				Description: "Creating route index file",
				Content: `import express from "express";
import roleRoutes from "@/app/role/role.routes";
import usersRoutes from "@/app/users/users.routes";
//IMPORT_AREA

const router = express.Router();
interface IRoute {
    path: string;
    route: any;
}
const moduleRoutes: IRoute[] = [
    {
        path: "/role",
        route: roleRoutes,
    },
    {
        path: "/users",
        route: usersRoutes,
    },//REGISTER_PATH_AREA
];

moduleRoutes.forEach((route) => router.use(route.path, route.route));

export default router;
`,
			},
}
 var updateInstructions=[]types.UpdateInstruction{

 }
 return fileInstructions,updateInstructions
}