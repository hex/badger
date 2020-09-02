using System;
using System.CommandLine;
using System.CommandLine.Invocation;
using System.IO;
using System.Linq;
using ImageMagick;

namespace Badger
{
    internal class BadgeOptions
    {
        public string Text { get; set; }
        public string Icon { get; set; }
        public string Color { get; set; }
        public string TextColor { get; set; }
        private string _textAlignment;

        public string TextAlignment
        {
            get => _textAlignment;
            set =>
                _textAlignment = value switch
                {
                    "top" => "North",
                    "left" => "West",
                    "bottom" => "South",
                    "right" => "East",
                    "topLeft" => "Northwest",
                    "topRight" => "Northeast",
                    "bottomLeft" => "Southwest",
                    "bottomRight" => "Southeast",
                    "center" => "Center",
                    _ => _textAlignment
                };
        }

        public int Angle { get; set; }
        public int OffsetX { get; set; }
        public int OffsetY { get; set; }

        private string _position;

        public string Position
        {
            get => _position;
            set =>
                _position = value switch
                {
                    "top" => "North",
                    "left" => "West",
                    "bottom" => "South",
                    "right" => "East",
                    "topLeft" => "Northwest",
                    "topRight" => "Northeast",
                    "bottomLeft" => "Southwest",
                    "bottomRight" => "Southeast",
                    "center" => "Center",
                    _ => _position
                };
        }

        public bool Replace { get; set; }
    }

    internal static class Program
    {
        private const string OutputDir = "badgerOutput";

        private static void AddAllOptions(Command command)
        {
            command.AddOption(Color());
            command.AddOption(TextColor());
            command.AddOption(TextAlignment());
            command.AddOption(Angle());
            command.AddOption(OffsetX());
            command.AddOption(OffsetY());
            command.AddOption(Position());
            command.AddOption(Replace());

            Option Color() =>
                new Option<string>(new[] {"-c", "--color"}, () => "#4096EE",
                    "Set badge color with a hexadecimal color code");

            Option TextColor() => new Option<string>(new[] {"-t", "--text-color",},
                () => "#F9F7ED",
                "Set badge text color with a hexadecimal color code");

            Option TextAlignment() => new Option<string>(new[] {"-l", "--text-alignment",},
                () => "center",
                "Set badge text alignment");

            Option Angle() => new Option<int>(new[] {"-a", "--angle"}, () => 0, "Set badge rotation");
            Option OffsetX() => new Option<int>(new[] {"-x", "--offset-x"}, () => 0, "Set badge x-axis offset");
            Option OffsetY() => new Option<int>(new[] {"-y", "--offset-y"}, () => 0, "Set badge y-axis offset");
            Option Position() => new Option<string>(new[] {"-p", "--position"}, () => "bottom", "Set badge position");
            Option Replace() => new Option<bool>(new[] {"-r", "--replace"}, () => false, "Replace input icon");
        }

        private static void EmptyDirectory(string path)
        {
            var di = new DirectoryInfo(path);

            foreach (var file in di.GetFiles())
            {
                file.Delete();
            }

            foreach (var dir in di.GetDirectories())
            {
                dir.Delete(true);
            }
        }

        private static void CreateBadge(BadgeOptions options)
        {
            Directory.CreateDirectory(OutputDir);
            EmptyDirectory(OutputDir);

            var topSettings = new MagickReadSettings
            {
                BackgroundColor = new MagickColor(options.Color),
                FontPointsize = 40,
                FillColor = new MagickColor(options.Color),
                Width = 1520
            };

            var bottomSettings = new MagickReadSettings
            {
                BackgroundColor = new MagickColor(options.Color),
                TextGravity = Enum.Parse<Gravity>(options.TextAlignment),
                FontWeight = FontWeight.Bold,
                FontPointsize = 180,
                AntiAlias = true,
                FillColor = new MagickColor(options.TextColor),
                TextInterlineSpacing = 10,
                Width = 1520
            };

            using (var top = new MagickImage($"caption:-", topSettings))
            {
                using (var bottom = new MagickImage($"caption:{options.Text}", bottomSettings))
                {
                    using (var badge = new MagickImageCollection())
                    {
                        bottom.BorderColor = new MagickColor(options.Color);
                        bottom.Border(100, 0);

                        badge.Add(top);
                        badge.Add(bottom);

                        using (var result = badge.AppendVertically())
                        {
                            result.BackgroundColor = MagickColors.Transparent;
                            // result.Rotate(options.Angle);

                            result.Write($"{OutputDir}{Path.DirectorySeparatorChar}badge.png");
                        }
                    }
                }
            }
        }

        private static void AppendBadge(BadgeOptions options, string path)
        {
            using (var icon = new MagickImage(path))
            {
                var badge = new MagickImage($"{OutputDir}{Path.DirectorySeparatorChar}badge.png");

                badge.Resize(icon.Width + icon.Width / 2, icon.Height + icon.Height / 2);
                badge.BackgroundColor = MagickColors.Transparent;
                badge.Rotate(options.Angle);

                var offsetX = icon.Width * options.OffsetX / 100d;
                var offsetY = icon.Height * options.OffsetY / 100d;
                icon.Composite(
                    badge,
                    Enum.Parse<Gravity>(options.Position),
                    new PointD(offsetX, offsetY),
                    CompositeOperator.Over
                );

                Console.WriteLine(
                    $"Writing to: {(options.Replace ? path : $"{OutputDir}{Path.DirectorySeparatorChar}{Path.GetFileName(path)}")}");

                icon.Write(options.Replace
                    ? path
                    : $"{OutputDir}{Path.DirectorySeparatorChar}{Path.GetFileName(path)}");
            }
        }

        private static int Main(string[] args)
        {
            var rootCommand = new RootCommand
            {
                Description = "A command-line tool that adds labels to your app icon",
                Name = "badger"
            };

            rootCommand.AddArgument(new Argument<string>("text", "Set label text"));
            rootCommand.AddArgument(new Argument<string>("icon",
                "Set path to icon with format .png | .jpg | .jpeg | .appiconset"));

            AddAllOptions(rootCommand);

            rootCommand.Handler = CommandHandler.Create(
                (BadgeOptions badgeOptions) =>
                {
                    CreateBadge(badgeOptions);

                    var fileAttributes = File.GetAttributes(badgeOptions.Icon);

                    if (fileAttributes.HasFlag(FileAttributes.Directory))
                    {
                        var files = Directory.EnumerateFiles(badgeOptions.Icon, "*.*", SearchOption.AllDirectories)
                            .Where(s => s.EndsWith(".png") || s.EndsWith(".jpg") || s.EndsWith(".jpeg"));

                        foreach (var file in files)
                        {
                            AppendBadge(badgeOptions, file);
                        }
                    }
                    else
                    {
                        AppendBadge(badgeOptions, badgeOptions.Icon);
                    }

                    if (badgeOptions.Replace)
                    {
                        Directory.Delete(OutputDir, true);
                    }
                    else
                    {
                        File.Delete($"{OutputDir}{Path.DirectorySeparatorChar}badge.png");
                    }
                });

            return rootCommand.InvokeAsync(args).Result;
        }
    }
}