'use client'
import { Button } from "@heroui/button";
import {Form} from "@heroui/form";
import { Input } from "@heroui/input";
import coverImg from "../../assets/signin.png"

import React from "react";


export default function Login() {
    const [action, setAction] = React.useState(null);
  return (
      <div className="container mx-auto">
        {/* wrapper */}
        <div className="flex flex-col lg:flex-row w-10/12 lg:w-8/12 bg-black-300 rounded-xl mx-auto shadow-lg overflow-hidden">
          {/* left */}
        
          <div className="w-full lg:w-1/2 py-16 px-12">
            <h2 className="text-3xl mb-4">Welcome Back</h2>
            <p className="mb-4">
             Log in and get on the vibe
            </p>
            <Form
      className="w-full max-w-xs flex flex-col gap-4"
      onReset={() => setAction("reset")}
      onSubmit={(e) => {
        e.preventDefault();
        let data = Object.fromEntries(new FormData(e.currentTarget));

        setAction(`submit ${JSON.stringify(data)}`);
      }}
    >
      

      <Input
        isRequired
        errorMessage="Please enter a valid email"
        label="Email"
        labelPlacement="outside"
        name="email"
        placeholder="Enter your email"
        type="email"
      />
      <Input
        isRequired
        errorMessage="Please enter your correct credentials"
        label="Password"
        labelPlacement="outside"
        name="password"
        placeholder="Enter your password"
        type="password"
      />
      <div className="flex gap-2">
        <Button color="danger" variant="bordered" size="lg" type="submit">
          Log In
        </Button>
        <Button type="reset" size="lg" variant="flat">
          Reset
        </Button>
      </div>
      {action && (
        <div className="text-small text-default-500">
          Action: <code>{action}</code>
        </div>
      )}
    </Form>
          </div>
  {/* right */}
           <div
            className="w-full lg:w-1/2 flex flex-col items-center justify-center p-12 bg-no-repeat bg-center bg-cover"
            style={{
              backgroundImage: `url(${coverImg.src})`,
            }}
          >
         <h1 className="text-white text-5xl font-bold mb-3 text-shadow-lg">
  Welcome
</h1>

          </div>
        </div>
      </div>
  );
}
