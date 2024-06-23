import tkinter as tk
import requests
from tkinter import scrolledtext

def on_button_click():
    chat_text.delete(1.0, tk.END)
    
    url = "http://localhost:5000/history"
    try:
        response = requests.get(url)
        if response.status_code == 200:
            history = response.text
            chat_text.insert(tk.END, history)
        else:
            chat_text.insert(tk.END, f"無法獲取歷史記錄。狀態碼: {response.status_code}\n")
    except requests.exceptions.RequestException as e:
        chat_text.insert(tk.END, f"發生錯誤：{str(e)}\n")

root = tk.Tk()
root.title("websocket 聊天室")
root.geometry("600x800")

chat_text = scrolledtext.ScrolledText(root, wrap=tk.WORD)
chat_text.pack(expand=True, fill=tk.BOTH, padx=10, pady=10)

button = tk.Button(root, text="Dump", command=on_button_click)
button.pack(anchor='ne', padx=10, pady=10)

root.mainloop()