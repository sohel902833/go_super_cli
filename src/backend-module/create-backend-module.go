package backendmodule

import (
	"github.com/sohel902833/go_super_cli/src/types"
)

func GetCreateBackendModuleInstructions()([]types.FileInstruction,[]types.UpdateInstruction){
  createInstructions:= []types.FileInstruction{
	 {
		FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.controller.ts",
		Description: "Create controller file",
		Content: `import { Request, Response, NextFunction } from "express";
import {
    create{{PASCAL_CASE_MODULE_NAME}}DTOSchema,
    edit{{PASCAL_CASE_MODULE_NAME}}DTOSchema,
} from "./{{LOWER_CASE_MODULE_NAME}}.schema";
import * as {{CAMEL_CASE_MODULE_NAME}}Service from "./{{LOWER_CASE_MODULE_NAME}}.service";
export const createNew{{PASCAL_CASE_MODULE_NAME}} = async (
    req: Request,
    res: Response,
    next: NextFunction
): Promise<any> => {
    try {
        const parsedBody = create{{PASCAL_CASE_MODULE_NAME}}DTOSchema.safeParse(req.body);
        if (!parsedBody.success) {
            return next(parsedBody.error);
        }
        const new{{PASCAL_CASE_MODULE_NAME}} = {
            ...parsedBody.data,
            creator: req.userId as string,
        };
        //@ts-ignore
        const created{{PASCAL_CASE_MODULE_NAME}} = await {{CAMEL_CASE_MODULE_NAME}}Service.create(new{{PASCAL_CASE_MODULE_NAME}});

        return res.status(201).json({
            message: "{{PASCAL_CASE_MODULE_NAME}} Successfully Created",
            data: created{{PASCAL_CASE_MODULE_NAME}},
            success: true,
        });
    } catch (err) {
        next(err);
    }
};

export const update{{PASCAL_CASE_MODULE_NAME}} = async (
    req: Request,
    res: Response,
    next: NextFunction
): Promise<any> => {
    try {
        const id = req.params.id as string;
        const parsedBody = edit{{PASCAL_CASE_MODULE_NAME}}DTOSchema.safeParse(req.body);
        if (!parsedBody.success) {
            return next(parsedBody.error);
        }
        //@ts-ignore
        const updated{{PASCAL_CASE_MODULE_NAME}} = await {{CAMEL_CASE_MODULE_NAME}}Service.edit(id, parsedBody.data);

        return res.json({
            message: "{{PASCAL_CASE_MODULE_NAME}} Successfully Updated",
            data: updated{{PASCAL_CASE_MODULE_NAME}},
            success: true,
        });
    } catch (err) {
        next(err);
    }
};

export const delete{{PASCAL_CASE_MODULE_NAME}} = async (
    req: Request,
    res: Response,
    next: NextFunction
): Promise<any> => {
    try {
        const id = req.params.id as string;
        const deleted{{PASCAL_CASE_MODULE_NAME}} = await {{CAMEL_CASE_MODULE_NAME}}Service.delete{{PASCAL_CASE_MODULE_NAME}}(id);
        return res.json({
            message: "{{PASCAL_CASE_MODULE_NAME}} Successfully Deleted",
            data: deleted{{PASCAL_CASE_MODULE_NAME}},
        });
    } catch (err) {
        next(err);
    }
};

export const getSingle{{PASCAL_CASE_MODULE_NAME}} = async (
    req: Request,
    res: Response,
    next: NextFunction
): Promise<any> => {
    try {
        const id = req.params.id as string;
        const {{PASCAL_CASE_MODULE_NAME}} = await {{CAMEL_CASE_MODULE_NAME}}Service.getSingle(id);
        return res.json({{PASCAL_CASE_MODULE_NAME}});
    } catch (err) {
        next(err);
    }
};

export const getAll{{PASCAL_CASE_MODULE_NAME}} = async (
    req: Request,
    res: Response,
    next: NextFunction
): Promise<any> => {
    try {
        const query = req.query;
        const {{PASCAL_CASE_MODULE_NAME}}s = await {{CAMEL_CASE_MODULE_NAME}}Service.getAll(query);
        return res.json({{PASCAL_CASE_MODULE_NAME}}s);
    } catch (err) {
        next(err);
    }
};
`,
	 },
	 {
		FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.model.ts",
		Description: "Creating model file",
		Content:`import { model, Schema, Types } from "mongoose";
import { I{{PASCAL_CASE_MODULE_NAME}}, I{{PASCAL_CASE_MODULE_NAME}} } from "./{{LOWER_CASE_MODULE_NAME}}.types.ts";
import { MODEL_NAMES } from "@/db";

const {{PASCAL_CASE_MODULE_NAME}}Schema = new Schema<I{{CAMEL_CASE_MODULE_NAME}}>(
    {
      {{MONGOOSE_SCHEMA_FIELDS}}
    },
    {
        timestamps: true,
    }
);
const {{PASCAL_CASE_MODULE_NAME}}Model = model<{{PASCAL_CASE_MODULE_NAME}}, I{{PASCAL_CASE_MODULE_NAME}}Model>(
    MODEL_NAMES.{{UPPER_CASE_MODULE_NAME}},
    {{PASCAL_CASE_MODULE_NAME}}Schema
);
export default {{PASCAL_CASE_MODULE_NAME}}Model;`,
	 },
	 {
		FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.routes.ts",
		Description: "Creating routes file",
		Content: `import express from "express";
import * as {{CAMEL_CASE_MODULE_NAME}}Controller from "./{{LOWER_CASE_MODULE_NAME}}.controller";
import { authGard } from "@/middlewares/authGard";
import { Permissions } from "../role";
const router = express.Router();
router.post(
    "/",
    authGard([Permissions.CREATE_{{UPPER_CASE_MODULE_NAME}}]),
    {{CAMEL_CASE_MODULE_NAME}}Controller.createNew{{PASCAL_CASE_MODULE_NAME}}
);
router.put(
    "/:id",
    authGard([Permissions.UPDATE_{{UPPER_CASE_MODULE_NAME}}]),
    {{CAMEL_CASE_MODULE_NAME}}Controller.update{{PASCAL_CASE_MODULE_NAME}}
);
router.delete(
    "/:id",
    authGard([Permissions.DELETE_{{UPPER_CASE_MODULE_NAME}}]),
    {{CAMEL_CASE_MODULE_NAME}}Controller.delete{{PASCAL_CASE_MODULE_NAME}}
);
router.get("/single/:id", {{CAMEL_CASE_MODULE_NAME}}Controller.getSingle{{PASCAL_CASE_MODULE_NAME}});
router.get("/", {{CAMEL_CASE_MODULE_NAME}}Controller.getAll{{PASCAL_CASE_MODULE_NAME}});

export default router;
`,
	 },
	 {
		FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.schema.ts",
		Description: "Creating schema file",
		Content: `import { z } from 'zod';

{{ZOD_GENERATED_SCHEMA}}

{{ZOD_INFER_TYPES}}

{{ZOD_EXPORTS}}`,
	 },
	 {
		FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.service.ts",
		Description: "Creating service file",
		Content: `import * as db from "@/db";
import { I{{PASCAL_CASE_MODULE_NAME}}, I{{PASCAL_CASE_MODULE_NAME}} } from "{{LOWER_CASE_MODULE_NAME}}.types";
import { QueryOptions } from "mongoose";
import { modifyQuery } from "@/helpers";
import { getWithPagination } from "@/helpers/pagination";

export const create = async (payload: I{{PASCAL_CASE_MODULE_NAME}}) => {
    try {
        const created{{PASCAL_CASE_MODULE_NAME}} = await db.models.{{PASCAL_CASE_MODULE_NAME}}Model.create(payload);
        return created{{PASCAL_CASE_MODULE_NAME}};
    } catch (err) {
        throw err;
    }
};

export const edit = async (id: string, payload: IEdit{{PASCAL_CASE_MODULE_NAME}}) => {
    try {
        const updated{{PASCAL_CASE_MODULE_NAME}} = await db.models.{{PASCAL_CASE_MODULE_NAME}}Model.findByIdAndUpdate(
            id,
            {
                $set: payload,
            },
            {
                new: true,
            }
        );
        return updated{{PASCAL_CASE_MODULE_NAME}};
    } catch (err) {
        throw err;
    }
};

export const delete{{PASCAL_CASE_MODULE_NAME}} = async (id: string) => {
    try {
        const deleted{{PASCAL_CASE_MODULE_NAME}} = await db.models.{{PASCAL_CASE_MODULE_NAME}}Model.findByIdAndDelete(id);
        return deleted{{PASCAL_CASE_MODULE_NAME}};
    } catch (err) {
        throw err;
    }
};

export const getSingle = async (id: string) => {
    try {
        const {{PASCAL_CASE_MODULE_NAME}} = await db.models.{{PASCAL_CASE_MODULE_NAME}}Model.findById(id).populate([
            {
                path: "creator",
                select: "firstName lastName",
            },
            {
                path: "category",
                select: "title",
            },
            {
                path: "subCategory",
                select: "title",
            },
            {
                path: "delivery_prices",
            },
        ]);
        return {{PASCAL_CASE_MODULE_NAME}};
    } catch (err) {
        throw err;
    }
};

export const getAll = async (filter: QueryOptions<I{{PASCAL_CASE_MODULE_NAME}}>) => {
    try {
        const { page, limit, finalQuery, sort } = modifyQuery(filter);
        const res = await getWithPagination({
            page: page,
            limit: limit,
            filter: finalQuery,
            model: db.models.{{PASCAL_CASE_MODULE_NAME}}Model,
            sort: sort,
            populate: [
            ],
        });
        return res;
    } catch (err) {
        throw err;
    }
};
`,
	 },
	 {
		FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.types.ts",
		Description: "Creating types file",
		Content: ``,
	 },
  }
  updateInstructions:=[]types.UpdateInstruction{};
  return createInstructions,updateInstructions
}