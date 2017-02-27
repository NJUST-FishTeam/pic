# coding: utf-8

from __future__ import print_function

import imghdr
import os
import hashlib
import sys

import tornado.web
import tornado.httpserver
import tornado.ioloop
import tornado.options
import tornado.autoreload

tornado.options.define("port", default=8888, help="run on the given port", type=int)
tornado.options.define("path", default="img", help="path to save img", type=str)
tornado.options.define("static", default="/static/images/", help="static ufl prefix", type=str)


class UploadHandler(tornado.web.RequestHandler):

    def post(self, *args, **kwargs):
        f = self.request.files['file'][0]
        data = f.get('body', '')
        img_type = imghdr.what('', h=data)
        if img_type:
            md5 = hashlib.md5(data).hexdigest()
            first, second, name = md5[0:2], md5[2:4], md5[4:] + '.' + img_type
            fpath = os.path.join(tornado.options.options.path, first, second)
            if not os.path.exists(fpath):
                os.makedirs(fpath)
            with open(os.path.join(fpath, name), 'w') as fp:
                fp.write(data)
            self.write(tornado.options.options.static + os.path.join(first, second, name))
        else:
            # 不是jpeg, bmp, png, rgb, gif, pbm, pgm, ppm, tiff, rast, xbm
            self.write(u'文件上传失败：不支持的文件类型')


class Application(tornado.web.Application):

    def __init__(self):
        settings = dict(
            template_path=os.path.join(os.path.dirname(__file__), "templates"),
            # static_path=os.path.join(os.path.dirname(__file__), "static"),
            debug=True
        )
        handlers = [
            (r'/upload/', UploadHandler),
            (r'/(.*)', tornado.web.StaticFileHandler, {
                'path': os.getcwd(),
                'default_filename': 'index.html',
            })
        ]
        tornado.web.Application.__init__(self, handlers, **settings)


if __name__ == '__main__':
    tornado.options.parse_command_line()
    if not os.path.exists(tornado.options.options.path):
        print('no such directory')
        sys.exit(0)
    http_server = tornado.httpserver.HTTPServer(Application())
    http_server.listen(tornado.options.options.port)
    instance = tornado.ioloop.IOLoop.instance()
    tornado.autoreload.start(instance)
    instance.start()